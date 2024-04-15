package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// ResetPasswordRequestBody is the request body for the ResetPassword route
type ResetPasswordRequestBody struct {
	Password string `json:"password" required:"true"`
}

// ResetPassword resets a user's password
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
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
