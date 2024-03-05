package authentication

import (
	"fmt"

	"github.com/gin-gonic/gin"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
)

func RegisterRoutes(router *gin.Engine, config *commonConfig.Config) (*ServiceClient, error) {
	client, err := InitServiceClient(config)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize authentication service client: %v", err)
	}
	service := &ServiceClient{
		client: client,
	}

	userRoutes := router.Group("/user")

	userRoutes.POST("/", service.Register)
	userRoutes.GET("/email/:verification_token", service.VerifyEmail)
	userRoutes.POST("/authenticate", service.Authenticate)

	emailRoutes := router.Group("/email")
	authenticationMiddleware, err := InitAuthenticationMiddleware(service)
	if err != nil {
		return nil, fmt.Errorf("Failed to initiate authenticator middleware: %v", err)
	}
	emailRoutes.Use(authenticationMiddleware.RequireAuthentication)
	emailRoutes.POST("/verification", service.ResendEmailVerification)

	authenticationRoutes := router.Group("/authentication")
	authenticationRoutes.Use(authenticationMiddleware.RefreshAuthentication)
	authenticationRoutes.POST("/refresh", service.RefreshToken)

	return service, nil
}
