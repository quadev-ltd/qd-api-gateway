package image_analysis

import (
	"context"
	"fmt"

	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_image_analysis"
	commonTLS "github.com/quadev-ltd/qd-common/pkg/tls"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/shared/middleware"
)

type ServiceClienter interface {
	ProcessImageAndPrompt(ctx context.Context, firebaseToken string, imageData []byte, prompt string) (*pb_image_analysis.ImagePromptResponse, error)
	middleware.ServiceClienter
}

type ServiceClient struct {
	client pb_image_analysis.ImageAnalysisServiceClient
}

var _ ServiceClienter = &ServiceClient{}

func InitServiceClient(configurations *config.Config) (ServiceClienter, error) {
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

func (service *ServiceClient) ProcessImageAndPrompt(ctx context.Context, firebaseToken string, imageData []byte, prompt string) (*pb_image_analysis.ImagePromptResponse, error) {
	return service.client.ProcessImageAndPrompt(ctx, &pb_image_analysis.ImagePromptRequest{
		FirebaseToken: firebaseToken,
		ImageData:     imageData,
		Prompt:        prompt,
	})
}

func (service *ServiceClient) GetPublicKey(ctx context.Context) (*string, error) {
	return nil, fmt.Errorf("image analysis service does not provide public keys")
}
