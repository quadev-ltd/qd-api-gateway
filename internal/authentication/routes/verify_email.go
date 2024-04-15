package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// VerifyEmail verifies an email
func VerifyEmail(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	res, err := client.VerifyEmail(
		ctx.Request.Context(),
		&pb_authentication.VerifyEmailRequest{
			UserId:            ctx.Param("user_id"),
			VerificationToken: ctx.Param("verification_token"),
		},
	)

	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
