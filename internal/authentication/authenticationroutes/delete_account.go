package authenticationroutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// DeleteAccount updates a user's profile
func DeleteAccount(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {

	res, err := client.DeleteAccount(
		ctx.Request.Context(),
		&pb_authentication.DeleteAccountRequest{},
	)

	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
