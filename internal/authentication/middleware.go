package authentication

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	commonJWT "github.com/quadev-ltd/qd-common/pkg/jwt"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"
)

type AutheticationMiddlewarer interface {
	RequireAuthentication(ctx *gin.Context)
	RefreshAuthentication(ctx *gin.Context)
}

type AutheticationMiddleware struct {
	service           ServiceClienter
	jwtVerifier       commonJWT.TokenVerifierer
	jwtTokenInspector commonJWT.TokenInspectorer
}

var _ AutheticationMiddlewarer = &AutheticationMiddleware{}

func InitAuthenticationMiddleware(authenticationService ServiceClienter) (AutheticationMiddlewarer, error) {
	correlationID := uuid.New().String()
	publicKey, err := RequestPublicKey(authenticationService, correlationID)
	if err != nil {
		return nil, err
	}
	jwtVerifier, err := commonJWT.NewTokenVerifier(*publicKey)
	if err != nil {
		return nil, err
	}
	jwtTokenInspector := &commonJWT.TokenInspector{}
	return &AutheticationMiddleware{
		authenticationService,
		jwtVerifier,
		jwtTokenInspector,
	}, nil
}

func RequestPublicKey(service ServiceClienter, correlationID string) (*string, error) {
	ctx := commonLogger.AddCorrelationIDToOutgoingContext(context.Background(), correlationID)
	publicKey, err := service.GetPublicKey(ctx)

	if err != nil {
		return nil, fmt.Errorf("Could not obtain public key: %v", err)
	}

	return publicKey, nil
}

func (autheticationMiddleware *AutheticationMiddleware) RequireAuthentication(ctx *gin.Context) {
	autheticationMiddleware.verifyToken(ctx, commonJWT.AccessTokenType)
}

func (autheticationMiddleware *AutheticationMiddleware) RefreshAuthentication(ctx *gin.Context) {
	autheticationMiddleware.verifyToken(ctx, commonJWT.RefreshTokenType)
}

func (autheticationMiddleware *AutheticationMiddleware) verifyToken(ctx *gin.Context, expectedTokenType commonJWT.TokenType) {
	logger, err := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	authorization := ctx.Request.Header.Get("Authorization")

	if authorization == "" {
		logger.Error(nil, "No authorization header was present in the request")
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		logger.Error(nil, "No bearer token was present in the authorization header")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	parsedToken, err := autheticationMiddleware.jwtVerifier.Verify(token[1])
	if err != nil {
		logger.Error(err, "The bearer token was invalid")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenType, err := autheticationMiddleware.jwtTokenInspector.GetTypeFromToken(parsedToken)
	if err != nil {
		logger.Error(err, "Could not obtain type from bearer token")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if commonJWT.TokenType(*tokenType) != expectedTokenType {
		logger.Error(nil, fmt.Sprintf("The bearer token was not an %s", expectedTokenType))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if expectedTokenType == commonJWT.RefreshTokenType {
		ctx.Set("token", token[1])
	}
	userEmail, err := autheticationMiddleware.jwtTokenInspector.GetEmailFromToken(parsedToken)
	if err != nil {
		logger.Error(err, "Could not obtain email from bearer token")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenExpiry, err := autheticationMiddleware.jwtTokenInspector.GetExpiryFromToken(parsedToken)
	if err != nil {
		logger.Error(err, "Could not obtain expiry from bearer token")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if tokenExpiry.Before(time.Now()) {
		logger.Error(nil, "The bearer token has expired")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set("userEmail", userEmail)

	logger.Info("Successfully authenticated user")
	ctx.Next()
}
