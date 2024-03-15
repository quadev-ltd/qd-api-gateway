package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

type UpdateUserProfileRequestBody struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth *int64 `json:"date_of_birth,omitempty"`
}

func UpdateUserProfile(ctx *gin.Context, client pb_authentication.AuthenticationServiceClient) {
	body := UpdateUserProfileRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	dateOfBirth := time.Unix(*body.DateOfBirth, 0)
	dateOfBirthProto := timestamppb.New(dateOfBirth)
	res, err := client.UpdateUserProfile(
		ctx.Request.Context(),
		&pb_authentication.UpdateUserProfileRequest{
			FirstName:   body.FirstName,
			LastName:    body.LastName,
			DateOfBirth: dateOfBirthProto,
		},
	)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
