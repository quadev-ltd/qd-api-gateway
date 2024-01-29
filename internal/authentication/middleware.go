package authentication

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	commonJTW "github.com/quadev-ltd/qd-common/pkg/jwt"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

type AutheticationMiddlewarer interface {
	RequireAuthentication(ctx *gin.Context)
}

type AutheticationMiddleware struct {
	service          *ServiceClient
	jwtAuthenticator commonJTW.JWTAthenticatorer
}

var _ AutheticationMiddlewarer = &AutheticationMiddleware{}

func InitAuthenticationMiddleware(authenticationService *ServiceClient) (AutheticationMiddlewarer, error) {
	correlationID := uuid.New().String()
	publicKey, err := RequestPublicKey(authenticationService, correlationID)
	if err != nil {
		return nil, err
	}
	jwtAuthenticator, err := commonJTW.NewJWTAuthenticator(*publicKey)
	if err != nil {
		return nil, err
	}
	return &AutheticationMiddleware{
		authenticationService,
		jwtAuthenticator,
	}, nil
}

func RequestPublicKey(service *ServiceClient, correlationID string) (*string, error) {
	ctx := commonLogger.AddCorrelationIDToOutgoingContext(context.Background(), correlationID)
	res, err := service.Client.GetPublicKey(
		ctx,
		&pb_authentication.GetPublicKeyRequest{},
	)

	if err != nil {
		return nil, fmt.Errorf("Could not obtain public key: %v", err)
	}

	return &res.PublicKey, nil
}

func (autheticationMiddleware *AutheticationMiddleware) RequireAuthentication(ctx *gin.Context) {
	logger := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	if logger == nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
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

	parsedToken, err := autheticationMiddleware.jwtAuthenticator.VerifyToken(token[1])
	if err != nil {
		logger.Error(err, "The bearer token was invalid")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userEmail, err := autheticationMiddleware.jwtAuthenticator.GetEmailFromToken(parsedToken)
	if err != nil {
		logger.Error(err, "Could not obtain email from bearer token")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenExpiry, err := autheticationMiddleware.jwtAuthenticator.GetExpiryFromToken(parsedToken)
	if err != nil {
		logger.Error(err, "Could not obatain expiry from bearer token")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if tokenExpiry.Compare(time.Now()) > 1 {
		logger.Error(nil, "The bearer token has expired")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set("userEmail", userEmail)

	logger.Info("Successfully authenticated user")
	ctx.Next()
}
