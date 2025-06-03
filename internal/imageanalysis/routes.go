package imageanalysis

import (
	"github.com/gin-gonic/gin"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	"golang.org/x/time/rate"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/middleware"
)

// RegisterRoutes registers all image analysis related routes with the provided router group
func RegisterRoutes(service ServiceClienter, api *gin.RouterGroup, configurations *commonConfig.Config, authMiddleware middleware.AutheticationMiddlewarer) error {
	rl := middleware.NewRateLimiter(rate.Limit(0.05), 3) // 3 requests per minute

	imageAnalysisRoutes := api.Group("/image-analysis")
	imageAnalysisRoutes.POST("", authMiddleware.RequirePaidFeatures, middleware.RateLimitMiddleware(rl), service.ProcessImageAndPrompt)

	return nil
}
