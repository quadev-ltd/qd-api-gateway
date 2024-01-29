package authentication

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
)

func RegisterRoutes(router *gin.Engine, config *config.Config) (*ServiceClient, error) {
	client, err := InitServiceClient(config)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize authentication service client: %v", err)
	}
	service := &ServiceClient{
		Client: client,
	}

	routes := router.Group("/user")

	routes.POST("/", service.Register)
	routes.GET("/email/:verification_token", service.VerifyEmail)

	routes = router.Group("/email")
	authenticationMiddleware, err := InitAuthenticationMiddleware(service)
	if err != nil {
		return nil, fmt.Errorf("Failed to initiate authenticator middleware: %v", err)
	}
	routes.Use(authenticationMiddleware.RequireAuthentication)
	routes.POST("/verification", service.ResendEmailVerification)

	return service, nil
}
