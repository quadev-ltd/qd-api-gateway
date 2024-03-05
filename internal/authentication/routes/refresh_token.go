package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

type RefreshTokentBody struct {
	Token string `json:"token"`
}

func RefreshToken(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	fmt.Println("token: ", ctx.GetString("token"))

	res, err := client.RefreshToken(ctx.Request.Context(), &pb_authentication.RefreshTokenRequest{
		Token: ctx.GetString("token"),
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
