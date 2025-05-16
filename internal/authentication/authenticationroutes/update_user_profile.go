package authenticationroutes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// UpdateUserProfileRequestBody is the request body for the UpdateUserProfile route
type UpdateUserProfileRequestBody struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth *int64 `json:"dateOfBirth,omitempty"`
}

// UpdateUserProfile updates a user's profile
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
		errors.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
