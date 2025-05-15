package image_analysis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_image_analysis"
	commonTLS "github.com/quadev-ltd/qd-common/pkg/tls"
	"github.com/rs/zerolog/log"

	commonConfig "github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/image_analysis/routes"
)

type ServiceClienter interface {
	ProcessImageAndPrompt(ctx *gin.Context)
}

type ServiceClient struct {
	client pb_image_analysis.ImageAnalysisServiceClient
}

var _ ServiceClienter = &ServiceClient{}

func InitServiceClient(configurations *commonConfig.Config) (ServiceClienter, error) {
	log.Info().Msg("Initializing image analysis service client")

	addr := fmt.Sprintf("%s:%s",
		configurations.ImageAnalysisService.Host,
		configurations.ImageAnalysisService.Port)

	conn, err := commonTLS.CreateGRPCConnection(addr, configurations.TLSEnabled)
	if err != nil {
		return nil, fmt.Errorf("Failed to create gRPC connection: %v", err)
	}

	client := pb_image_analysis.NewImageAnalysisServiceClient(conn)

	return &ServiceClient{
		client: client,
	}, nil
}

func (service *ServiceClient) ProcessImageAndPrompt(ctx *gin.Context) {
	routes.ProcessImageAndPrompt(ctx, service.client)
}
