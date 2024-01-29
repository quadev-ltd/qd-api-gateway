package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

type RegisterRequestBody struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth *int64 `json:"date_of_birth,omitempty"`
}

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
	dateOfBirth := time.Unix(*body.DateOfBirth, 0)
	dateOfBirthProto := timestamppb.New(dateOfBirth)
	res, err := client.Register(ctx.Request.Context(), &pb_authentication.RegisterRequest{
		Email:       body.Email,
		Password:    body.Password,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		DateOfBirth: dateOfBirthProto,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
