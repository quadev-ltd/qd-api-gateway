package authentication

import (
	"github.com/gin-gonic/gin"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	"golang.org/x/time/rate"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/middleware"
)

// RegisterRoutes registers the authentication routes
func RegisterRoutes(
	service ServiceClienter,
	api *gin.RouterGroup,
	centralConfig *commonConfig.Config,
	configurations *config.Config,
	authenticationMiddleware middleware.AutheticationMiddlewarer,
) error {

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
	userRoutes.DELETE("", authenticationMiddleware.RequireAuthentication, service.DeleteAccount)

	authenticationRoutes := api.Group("/authentication")
	authenticationRoutes.Use(authenticationMiddleware.RefreshAuthentication)
	authenticationRoutes.POST("/refresh", service.RefreshToken)

	return nil
}
