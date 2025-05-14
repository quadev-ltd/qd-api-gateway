package middleware

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	commonJWT "github.com/quadev-ltd/qd-common/pkg/jwt"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"
	commonToken "github.com/quadev-ltd/qd-common/pkg/token"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
)

type ServiceClienter interface {
	GetPublicKey(ctx context.Context) (*string, error)
}

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

func InitAuthenticationMiddleware(authenticationService ServiceClienter, configurations *config.Config) (AutheticationMiddlewarer, error) {
	correlationID := uuid.New().String()
	publicKey, err := RequestPublicKey(authenticationService, correlationID, configurations.Environment, backoffDelay)
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

type BackoffStrategy func(attempt int) time.Duration

func backoffDelay(attempt int) time.Duration {
	const maxDelay = 30 * time.Second
	delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

func RequestPublicKey(
	service ServiceClienter,
	correlationID,
	environment string,
	backoff BackoffStrategy,
) (*string, error) {
	logFactory := commonLogger.NewLogFactory(environment)
	logger := logFactory.NewLogger()
	var publicKey *string
	var err error

	maxAttempts := 5
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx := commonLogger.AddCorrelationIDToOutgoingContext(context.Background(), correlationID)
		publicKey, err = service.GetPublicKey(ctx)

		if err == nil {
			return publicKey, nil
		}

		logger.Info(fmt.Sprintf("Attempt %d: could not obtain public key, error: %v\n", attempt, err))
		time.Sleep(backoff(attempt))
	}

	return nil, fmt.Errorf("Could not obtain public key after %d attempts: %v", maxAttempts, err)
}

func (autheticationMiddleware *AutheticationMiddleware) RequireAuthentication(ctx *gin.Context) {
	autheticationMiddleware.verifyToken(ctx, commonToken.AuthTokenType)
}

func (autheticationMiddleware *AutheticationMiddleware) RefreshAuthentication(ctx *gin.Context) {
	autheticationMiddleware.verifyToken(ctx, commonToken.RefreshTokenType)
}

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

func (autheticationMiddleware *AutheticationMiddleware) verifyToken(ctx *gin.Context, expectedTokenType commonToken.Type) {
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
	if commonToken.Type(claims.Type) != expectedTokenType {
		logger.Error(nil, fmt.Sprintf("The bearer token was not an %s but a %s", expectedTokenType, claims.Type))
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
