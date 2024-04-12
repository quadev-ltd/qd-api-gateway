package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

// AuthenticateRequestBody is the request body for the Authenticate route
type AuthenticateRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Authenticate authenticates a user
func Authenticate(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	body := AuthenticateRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := client.Authenticate(ctx.Request.Context(), &pb_authentication.AuthenticateRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		errorHTTPStatusCode := errors.GRPCErrorToHTTPStatus(err)
		ctx.JSON(errorHTTPStatusCode, gin.H{"error": err.Error()})
		ctx.AbortWithError(errorHTTPStatusCode, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
