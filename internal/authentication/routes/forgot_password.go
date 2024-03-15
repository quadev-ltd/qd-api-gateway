package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

type ForgotPasswordRequestBody struct {
	Email string `json:"email" required:"true"`
}

func ForgotPassword(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	body := ForgotPasswordRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := client.ForgotPassword(
		ctx.Request.Context(),
		&pb_authentication.ForgotPasswordRequest{
			Email: body.Email,
		},
	)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}