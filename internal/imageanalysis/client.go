package imageanalysis

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_image_analysis"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	commonTLS "github.com/quadev-ltd/qd-common/pkg/tls"
	"github.com/rs/zerolog/log"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// ServiceClienter defines the interface for the image analysis service client
type ServiceClienter interface {
	// ProcessImageAndPrompt processes an image with a given prompt
	ProcessImageAndPrompt(ctx *gin.Context)
}

// ServiceClient implements the ServiceClienter interface and handles communication with the image analysis service
type ServiceClient struct {
	client pb_image_analysis.ImageAnalysisServiceClient
}

var _ ServiceClienter = &ServiceClient{}

// InitServiceClient initializes a new image analysis service client with the provided configuration
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

// ProcessImagePromptRequestBody represents the expected request body for image processing
type ProcessImagePromptRequestBody struct {
	Image  []byte `json:"image" binding:"required"`  // Base64 encoded image data
	Prompt string `json:"prompt" binding:"required"` // Text prompt for image analysis
}

// ProcessImageAndPrompt handles the image processing request by forwarding it to the image analysis service
func ProcessImageAndPrompt(ctx *gin.Context, client pb_image_analysis.ImageAnalysisServiceClient) {
	body := ProcessImagePromptRequestBody{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := client.ProcessImageAndPrompt(
		ctx.Request.Context(),
		&pb_image_analysis.ImagePromptRequest{
			FirebaseToken: "",
			ImageData:     body.Image,
			Prompt:        body.Prompt,
		},
	)

	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// ProcessImageAndPrompt handles the HTTP request to process an image with a prompt
func (service *ServiceClient) ProcessImageAndPrompt(ctx *gin.Context) {
	ProcessImageAndPrompt(ctx, service.client)
}
