package authentication

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/shared/middleware"
)

// AutheticationMiddlewarer interface is used to verify JWT tokens
type AutheticationMiddlewarer interface {
	middleware.AutheticationMiddlewarer
}

// ServiceClienter interface for the authentication service
type ServiceClienter interface {
	GetPublicKey(ctx context.Context) (*string, error)
	middleware.ServiceClienter
}

// InitAuthenticationMiddleware initializes the authentication middleware
func InitAuthenticationMiddleware(authenticationService ServiceClienter, configurations *config.Config) (AutheticationMiddlewarer, error) {
	return middleware.InitAuthenticationMiddleware(authenticationService, configurations)
}

// ParseAccessToken parses the access token from the request
func ParseAccessToken(ctx *gin.Context) *string {
	return middleware.ParseAccessToken(ctx)
}
