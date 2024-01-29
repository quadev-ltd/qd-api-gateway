package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

func ResendEmailVerification(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	authorization := ctx.Request.Header.Get("Authorization")
	authToken := strings.Split(authorization, "Bearer ")[1]

	res, err := client.ResendEmailVerification(ctx.Request.Context(), &pb_authentication.ResendEmailVerificationRequest{
		AuthToken: authToken,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
