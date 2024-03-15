package authentication

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	commonTLS "github.com/quadev-ltd/qd-common/pkg/tls"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/authentication/routes"
	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

type ServiceClienter interface {
	GetPublicKey(ctx context.Context) (*string, error)
	Register(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	ResendEmailVerification(ctx *gin.Context)
	Authenticate(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	VerifyResetPasswordToken(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	GetUserProfile(ctx *gin.Context)
	UpdateUserProfile(ctx *gin.Context)
}

type ServiceClient struct {
	client pb_authentication.AuthenticationServiceClient
}

var _ ServiceClienter = &ServiceClient{}

func InitServiceClient(config *commonConfig.Config) (pb_authentication.AuthenticationServiceClient, error) {
	grpcServiceAddress := fmt.Sprintf("%s:%s", config.AuthenticationService.Host, config.AuthenticationService.Port)

	clientConnection, err := commonTLS.CreateGRPCConnection(grpcServiceAddress, config.TLSEnabled)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to grpc authentication service: %v", err)
	}

	return pb_authentication.NewAuthenticationServiceClient(clientConnection), nil
}

func (service *ServiceClient) GetPublicKey(ctx context.Context) (*string, error) {
	response, err := service.client.GetPublicKey(
		ctx,
		&pb_authentication.GetPublicKeyRequest{},
	)
	if err != nil {
		return nil, err
	}
	return &response.PublicKey, nil
}

func (service *ServiceClient) Register(ctx *gin.Context) {
	routes.Register(ctx, service.client)
}

func (service *ServiceClient) VerifyEmail(ctx *gin.Context) {
	routes.VerifyEmail(ctx, service.client)
}

func (service *ServiceClient) ResendEmailVerification(ctx *gin.Context) {
	routes.ResendEmailVerification(ctx, service.client)
}

func (service *ServiceClient) Authenticate(ctx *gin.Context) {
	routes.Authenticate(ctx, service.client)
}

func (service *ServiceClient) RefreshToken(ctx *gin.Context) {
	routes.RefreshToken(ctx, service.client)
}

func (service *ServiceClient) ForgotPassword(ctx *gin.Context) {
	routes.ForgotPassword(ctx, service.client)
}

func (service *ServiceClient) VerifyResetPasswordToken(ctx *gin.Context) {
	routes.VerifyResetPasswordToken(ctx, service.client)
}

func (service *ServiceClient) ResetPassword(ctx *gin.Context) {
	routes.ResetPassword(ctx, service.client)
}

func (service *ServiceClient) GetUserProfile(ctx *gin.Context) {
	routes.GetUserProfile(ctx, service.client)
}

func (service *ServiceClient) UpdateUserProfile(ctx *gin.Context) {
	routes.UpdateUserProfile(ctx, service.client)
}
