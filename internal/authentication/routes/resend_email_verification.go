package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/util"
	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

func ResendEmailVerification(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	res, err := client.ResendEmailVerification(
		ctx.Request.Context(),
		&pb_authentication.ResendEmailVerificationRequest{},
	)

	if err != nil {
		errorHttpStatusCode := util.GRPCErrorToHTTPStatus(err)
		ctx.JSON(errorHttpStatusCode, gin.H{"error": err.Error()})
		ctx.AbortWithError(errorHttpStatusCode, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
