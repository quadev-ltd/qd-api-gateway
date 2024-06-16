package authentication

import (
	"fmt"

	"github.com/gin-gonic/gin"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	"golang.org/x/time/rate"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/middleware"
)

// RegisterRoutes registers the authentication routes
func RegisterRoutes(
	api *gin.RouterGroup,
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

	rl := middleware.NewRateLimiter(rate.Limit(0.08), 5)

	userRoutes := api.Group("/user")
	userRoutes.POST("/", middleware.RateLimitMiddleware(rl), service.Register)
	userRoutes.POST("/:userID/email/:verificationToken", service.VerifyEmail)
	userRoutes.POST("/sessions", middleware.RateLimitMiddleware(rl), service.Authenticate)
	userRoutes.POST("/firebase/sessions", middleware.RateLimitMiddleware(rl), service.AuthenticateWithFirebase)
	userRoutes.POST("/:userID/email/verification", middleware.RateLimitMiddleware(rl), service.ResendEmailVerification)
	userRoutes.POST("/password/reset", middleware.RateLimitMiddleware(rl), service.ForgotPassword)
	userRoutes.GET("/:userID/password/reset-verification/:verificationToken", middleware.RateLimitMiddleware(rl), service.VerifyResetPasswordToken)
	userRoutes.POST("/:userID/password/reset/:verificationToken", middleware.RateLimitMiddleware(rl), service.ResetPassword)
	userRoutes.GET("/profile", authenticationMiddleware.RequireAuthentication, service.GetUserProfile)
	userRoutes.PUT("/profile", authenticationMiddleware.RequireAuthentication, service.UpdateUserProfile)

	authenticationRoutes := api.Group("/authentication")
	authenticationRoutes.Use(authenticationMiddleware.RefreshAuthentication)
	authenticationRoutes.POST("/refresh", service.RefreshToken)

	return service, nil
}
