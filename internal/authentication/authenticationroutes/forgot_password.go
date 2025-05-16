package authenticationroutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// ForgotPasswordRequestBody is the request body for the ForgotPassword route
type ForgotPasswordRequestBody struct {
	Email string `json:"email" required:"true"`
}

// ForgotPassword requests a password reset email
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
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
