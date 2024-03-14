package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

func GetUserProfile(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	res, err := client.GetUserProfile(
		ctx.Request.Context(),
		&pb_authentication.GetUserProfileRequest{},
	)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
