package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// AuthenticateWithFirebaseRequestBody is the request body for the Authenticate route
type AuthenticateWithFirebaseRequestBody struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IDToken   string `json:"idToken"`
}

// AuthenticateWithFirebase authenticates a user using firebase
func AuthenticateWithFirebase(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	body := AuthenticateWithFirebaseRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := client.AuthenticateWithFirebase(
		ctx.Request.Context(),
		&pb_authentication.AuthenticateWithFirebaseRequest{
			Email:     body.Email,
			FirstName: body.FirstName,
			LastName:  body.LastName,
			IdToken:   body.IDToken,
		},
	)

	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
