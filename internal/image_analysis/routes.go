package image_analysis

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/middleware"
)

func RegisterRoutes(service ServiceClienter, api *gin.RouterGroup, configurations *config.Config, authMiddleware middleware.AutheticationMiddlewarer) error {
	rl := middleware.NewRateLimiter(rate.Limit(0.05), 3) // 3 requests per minute

	imageAnalysisRoutes := api.Group("/image-analysis")
	imageAnalysisRoutes.POST("", authMiddleware.RequireAuthentication, middleware.RateLimitMiddleware(rl), service.ProcessImageAndPrompt)

	return nil
}
