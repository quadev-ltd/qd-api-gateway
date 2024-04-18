package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// ResendEmailVerification resends an email verification
func ResendEmailVerification(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	res, err := client.ResendEmailVerification(
		ctx.Request.Context(),
		&pb_authentication.ResendEmailVerificationRequest{
			UserID: ctx.Param("userID"),
		},
	)

	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
