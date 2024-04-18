package authentication

import (
	"fmt"

	"github.com/gin-gonic/gin"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
)

// RegisterRoutes registers the authentication routes
func RegisterRoutes(
	router *gin.Engine,
	centralConfig *commonConfig.Config,
	configurations *config.Config,
) (*ServiceClient, error) {
	client, err := InitServiceClient(centralConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize authentication service client: %v", err)
	}
	service := &ServiceClient{
		client: client,
	}

	authenticationMiddleware, err := InitAuthenticationMiddleware(service, configurations)
	if err != nil {
		return nil, fmt.Errorf("Failed to initiate authenticator middleware: %v", err)
	}

	api := router.Group("/api/v1")

	userRoutes := api.Group("/user")
	userRoutes.POST("/", service.Register)
	userRoutes.POST("/:userID/email/:verificationToken", service.VerifyEmail)
	userRoutes.POST("/sessions", service.Authenticate)
	userRoutes.POST("/:userID/email/verification", service.ResendEmailVerification)
	userRoutes.POST("/password/reset-request", service.ForgotPassword)
	userRoutes.GET("/:userID/password/reset-verification/:verificationToken", service.VerifyResetPasswordToken)
	userRoutes.POST("/:userID/password/reset/:verificationToken", service.ResetPassword)
	userRoutes.GET("/", authenticationMiddleware.RequireAuthentication, service.GetUserProfile)
	userRoutes.PUT("/", authenticationMiddleware.RequireAuthentication, service.UpdateUserProfile)

	authenticationRoutes := api.Group("/authentication")
	authenticationRoutes.Use(authenticationMiddleware.RefreshAuthentication)
	authenticationRoutes.POST("/refresh", service.RefreshToken)

	return service, nil
}
