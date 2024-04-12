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
	commonToken "github.com/quadev-ltd/qd-common/pkg/token"
)

// AutheticationMiddlewarer interface is used to verify JWT tokens
type AutheticationMiddlewarer interface {
	RequireAuthentication(ctx *gin.Context)
	RefreshAuthentication(ctx *gin.Context)
}

// AutheticationMiddleware is used to verify JWT tokens
type AutheticationMiddleware struct {
	service           ServiceClienter
	jwtVerifier       commonJWT.TokenVerifierer
	jwtTokenInspector commonJWT.TokenInspectorer
}

var _ AutheticationMiddlewarer = &AutheticationMiddleware{}

// InitAuthenticationMiddleware initializes the authentication middleware
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

// RequestPublicKey requests the public key from the authentication service
func RequestPublicKey(service ServiceClienter, correlationID string) (*string, error) {
	ctx := commonLogger.AddCorrelationIDToOutgoingContext(context.Background(), correlationID)
	publicKey, err := service.GetPublicKey(ctx)

	if err != nil {
		return nil, fmt.Errorf("Could not obtain public key: %v", err)
	}

	return publicKey, nil
}

// RequireAuthentication verifies the access token
func (autheticationMiddleware *AutheticationMiddleware) RequireAuthentication(ctx *gin.Context) {
	autheticationMiddleware.verifyToken(ctx, commonToken.AccessTokenType)
}

// RefreshAuthentication verifies the refresh token
func (autheticationMiddleware *AutheticationMiddleware) RefreshAuthentication(ctx *gin.Context) {
	autheticationMiddleware.verifyToken(ctx, commonToken.RefreshTokenType)
}

// ParseAccessToken parses the access token from the request
func ParseAccessToken(ctx *gin.Context) *string {
	logger, err := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}
	authorization := ctx.Request.Header.Get("Authorization")

	if authorization == "" {
		logger.Error(nil, "No authorization header was present in the request")
		ctx.AbortWithError(
			http.StatusForbidden,
			fmt.Errorf("No authorization header was present in the request"),
		)
		return nil
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		logger.Error(nil, "No bearer token was present in the authorization header")
		ctx.AbortWithError(
			http.StatusUnauthorized,
			fmt.Errorf("No bearer token was present in the authorization header"),
		)
		return nil
	}
	return &token[1]
}

func (autheticationMiddleware *AutheticationMiddleware) verifyToken(ctx *gin.Context, expectedTokenType commonToken.TokenType) {
	logger, err := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	parsedAuthorizationToken := ParseAccessToken(ctx)
	if parsedAuthorizationToken == nil {
		return
	}
	parsedToken, err := autheticationMiddleware.jwtVerifier.Verify(*parsedAuthorizationToken)
	if err != nil {
		logger.Error(err, "The bearer token was invalid")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, err := autheticationMiddleware.jwtTokenInspector.GetClaimsFromToken(parsedToken)
	if err != nil {
		logger.Error(err, "Could not obtain claims from bearer token")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if commonToken.TokenType(claims.Type) != expectedTokenType {
		logger.Error(nil, fmt.Sprintf("The bearer token was not an %s", expectedTokenType))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims.Expiry.Before(time.Now()) {
		logger.Error(nil, "The bearer token has expired")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	newContext := commonJWT.AddAuthorizationMetadataToContext(ctx.Request.Context(), *parsedAuthorizationToken)
	ctx.Request = ctx.Request.WithContext(newContext)
	ctx.Set(string(commonJWT.ClaimsContextKey), claims)
	ctx.Set(string(commonJWT.JWTTokenKey), parsedToken)

	logger.Info("Successfully authenticated user")
	ctx.Next()
}
