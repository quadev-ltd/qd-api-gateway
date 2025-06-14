package middleware

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	commonJWT "github.com/quadev-ltd/qd-common/pkg/jwt"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"
	commonToken "github.com/quadev-ltd/qd-common/pkg/token"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
)

// ServiceClienter defines the interface for services that can provide public keys
type ServiceClienter interface {
	GetPublicKey(ctx context.Context) (*string, error)
}

// AutheticationMiddlewarer defines the interface for authentication middleware
type AutheticationMiddlewarer interface {
	// RequireAuthentication ensures the request has a valid authentication token
	RequireAuthentication(ctx *gin.Context)
	// RefreshAuthentication ensures the request has a valid refresh token
	RefreshAuthentication(ctx *gin.Context)
	// RequirePaidFeatures ensures the request has a valid authentication token with HasPaidFeatures set to true
	RequirePaidFeatures(ctx *gin.Context)
}

// AutheticationMiddleware implements the AutheticationMiddlewarer interface
type AutheticationMiddleware struct {
	service           ServiceClienter
	jwtVerifier       commonJWT.TokenVerifierer
	jwtTokenInspector commonJWT.TokenInspectorer
}

var _ AutheticationMiddlewarer = &AutheticationMiddleware{}

// InitAuthenticationMiddleware initializes a new authentication middleware with the provided service and configuration
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

// BackoffStrategy defines a function type for implementing backoff delays
type BackoffStrategy func(attempt int) time.Duration

// backoffDelay implements an exponential backoff strategy with a maximum delay
func backoffDelay(attempt int) time.Duration {
	const maxDelay = 30 * time.Second
	delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}

// RequestPublicKey attempts to retrieve the public key from the authentication service with retries
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

// RequireAuthentication middleware ensures the request has a valid authentication token
func (autheticationMiddleware *AutheticationMiddleware) RequireAuthentication(ctx *gin.Context) {
	_, isVerified := autheticationMiddleware.verifyToken(ctx, commonToken.AuthTokenType)
	if isVerified {
		ctx.Next()
	}
}

// RefreshAuthentication middleware ensures the request has a valid refresh token
func (autheticationMiddleware *AutheticationMiddleware) RefreshAuthentication(ctx *gin.Context) {
	_, isVerified := autheticationMiddleware.verifyToken(ctx, commonToken.RefreshTokenType)
	if isVerified {
		ctx.Next()
	}
}

// ParseAccessToken extracts and validates the access token from the request headers
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

func (autheticationMiddleware *AutheticationMiddleware) verifyToken(
	ctx *gin.Context,
	expectedTokenType commonToken.Type,
) (*jwt.Token, bool) {
	logger, err := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return nil, false
	}
	parsedAuthorizationToken := ParseAccessToken(ctx)
	if parsedAuthorizationToken == nil {
		return nil, false
	}
	parsedToken, err := autheticationMiddleware.jwtVerifier.Verify(*parsedAuthorizationToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "The bearer token was invalid")
		return nil, false
	}
	claims, err := autheticationMiddleware.jwtTokenInspector.GetClaimsFromToken(parsedToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Could not obtain claims from bearer token")
		return nil, false
	}
	if commonToken.Type(claims.Type) != expectedTokenType {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Sprintf("The bearer token was not an %s but a %s", expectedTokenType, claims.Type))
		return nil, false
	}

	if claims.Expiry.Before(time.Now()) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "The bearer token has expired")
		return nil, false
	}

	newContext := commonJWT.AddAuthorizationMetadataToContext(ctx.Request.Context(), *parsedAuthorizationToken)
	ctx.Request = ctx.Request.WithContext(newContext)
	ctx.Set(string(commonJWT.ClaimsContextKey), claims)
	ctx.Set(string(commonJWT.JWTTokenKey), parsedToken)

	logger.Info("Successfully authenticated user")
	return parsedToken, true
}

// RequirePaidFeatures middleware ensures the request has a valid authentication token with HasPaidFeatures set to true
func (autheticationMiddleware *AutheticationMiddleware) RequirePaidFeatures(ctx *gin.Context) {
	logger, err := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	parsedToken, isVerified := autheticationMiddleware.verifyToken(ctx, commonToken.AuthTokenType)
	if !isVerified {
		return
	}
	claims, err := autheticationMiddleware.jwtTokenInspector.GetClaimsFromToken(parsedToken)
	if err != nil {
		logger.Error(err, "Could not obtain claims from bearer token")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	hasPaidFeatures := bool(claims.HasPaidFeatures)
	if !hasPaidFeatures {
		ctx.AbortWithStatusJSON(http.StatusPaymentRequired, "User does not have access to paid features")
		return
	}

	logger.Info("User has access to paid features")
	ctx.Next()
}
