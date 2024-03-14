package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

func VerifyResetPasswordToken(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	res, err := client.VerifyResetPasswordToken(
		ctx.Request.Context(),
		&pb_authentication.VerifyResetPasswordTokenRequest{
			UserId: ctx.Param("user_id"),
			Token:  ctx.Param("verification_token"),
		},
	)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
