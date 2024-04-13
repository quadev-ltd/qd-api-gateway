package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// RegisterRequestBody is the request body for the Register route
type RegisterRequestBody struct {
	Email       string                 `json:"email"`
	Password    string                 `json:"password"`
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	DateOfBirth *timestamppb.Timestamp `json:"date_of_birth,omitempty"`
}

// Register registers a new user
func Register(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	body := RegisterRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if body.DateOfBirth == nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("date_of_birth is required"))
		return
	}
	res, err := client.Register(ctx.Request.Context(), &pb_authentication.RegisterRequest{
		Email:       body.Email,
		Password:    body.Password,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		DateOfBirth: body.DateOfBirth,
	})

	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
