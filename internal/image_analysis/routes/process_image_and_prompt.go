package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_image_analysis"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

type ProcessImagePromptRequestBody struct {
	Image  []byte `json:"image" binding:"required"`
	Prompt string `json:"prompt" binding:"required"`
}

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
