package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// GetUserProfile requests a user's profile
func GetUserProfile(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	res, err := client.GetUserProfile(
		ctx.Request.Context(),
		&pb_authentication.GetUserProfileRequest{},
	)

	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
