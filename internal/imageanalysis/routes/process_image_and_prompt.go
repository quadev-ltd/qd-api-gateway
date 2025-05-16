package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_image_analysis"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

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
