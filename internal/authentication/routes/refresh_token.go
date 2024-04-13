package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// RefreshTokentBody is the request body for the RefreshToken route
type RefreshTokentBody struct {
	Token string `json:"token"`
}

// RefreshToken refreshes a user's token
func RefreshToken(
	ctx *gin.Context,
	client pb_authentication.AuthenticationServiceClient,
) {
	res, err := client.RefreshToken(
		ctx.Request.Context(),
		&pb_authentication.RefreshTokenRequest{},
	)

	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
