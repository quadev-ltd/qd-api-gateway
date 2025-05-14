package image_analysis

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/image_analysis/routes"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/middleware"
	sharedMiddleware "github.com/quadev-ltd/qd-qpi-gateway/internal/shared/middleware"
)

type Service struct {
	ProcessImagePrompt gin.HandlerFunc
	client             ServiceClienter
}

func RegisterRoutes(api *gin.RouterGroup, configurations *config.Config, authMiddleware sharedMiddleware.AutheticationMiddlewarer) (*Service, error) {
	client, err := InitServiceClient(configurations)
	if err != nil {
		return nil, err
	}

	service := &Service{
		ProcessImagePrompt: func(ctx *gin.Context) {
			routes.ProcessImagePrompt(ctx, client)
		},
		client: client,
	}

	rl := middleware.NewRateLimiter(rate.Limit(0.05), 3) // 3 requests per minute

	imageAnalysisRoutes := api.Group("/image-analysis")
	imageAnalysisRoutes.POST("", authMiddleware.RequireAuthentication, middleware.RateLimitMiddleware(rl), service.ProcessImagePrompt)

	return service, nil
}
