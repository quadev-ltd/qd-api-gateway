package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

// GetUserProfile requests a user's profile
func GetUserProfile(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	res, err := client.GetUserProfile(
		ctx.Request.Context(),
		&pb_authentication.GetUserProfileRequest{},
	)

	if err != nil {
		errorHTTPStatusCode := errors.GRPCErrorToHTTPStatus(err)
		ctx.JSON(errorHTTPStatusCode, gin.H{"error": err.Error()})
		ctx.AbortWithError(errorHTTPStatusCode, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
