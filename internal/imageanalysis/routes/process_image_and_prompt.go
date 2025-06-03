package routes

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_image_analysis"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// ProcessImagePromptRequestBody represents the expected request body for image processing
type ProcessImagePromptRequestBody struct {
	Image  []byte `json:"image" binding:"required"`  // Base64 encoded image data
	Prompt string `json:"prompt" binding:"required"` // Text prompt for image analysis
}

// ProcessImageAndPrompt handles the image processing request by forwarding it to the image analysis service
func ProcessImageAndPrompt(ctx *gin.Context, client pb_image_analysis.ImageAnalysisServiceClient) {
	logger, err := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get the image file from the multipart form
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		logger.Error(err, "Error getting image file from multipart form")
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("File error"))
		return
	}
	defer file.Close()
	imageData, err := io.ReadAll(file)
	if err != nil {
		logger.Error(err, "Error reading image file")
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("File error"))
		return
	}
	prompt := ctx.PostForm("prompt")

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = ctx.PostForm("mimeType")
	}

	logger.Info(
		fmt.Sprintf(
			"Processing image size %d bytes, mimeType %s and prompt length %d characters",
			len(imageData),
			mimeType,
			len(prompt),
		),
	)
	res, err := client.ProcessImageAndPrompt(
		ctx.Request.Context(),
		&pb_image_analysis.ImagePromptRequest{
			ImageData: imageData,
			Prompt:    prompt,
			MimeType:  mimeType,
		},
	)

	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
