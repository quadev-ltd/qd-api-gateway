package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

type AuthenticateRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
