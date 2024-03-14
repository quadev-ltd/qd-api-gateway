package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

type ResetPasswordRequestBody struct {
	Password string `json:"password" required:"true"`
}

func ResetPassword(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	body := ResetPasswordRequestBody{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := client.ResetPassword(
		ctx.Request.Context(),
		&pb_authentication.ResetPasswordRequest{
			UserId:      ctx.Param("user_id"),
			Token:       ctx.Param("verification_token"),
			NewPassword: body.Password,
		},
	)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
