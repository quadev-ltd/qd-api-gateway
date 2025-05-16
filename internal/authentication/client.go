package authentication

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	commonTLS "github.com/quadev-ltd/qd-common/pkg/tls"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/authentication/authenticationroutes"
)

// ServiceClienter is an interface for the authentication service client
type ServiceClienter interface {
	GetPublicKey(ctx context.Context) (*string, error)
	Register(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	ResendEmailVerification(ctx *gin.Context)
	Authenticate(ctx *gin.Context)
	AuthenticateWithFirebase(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	VerifyResetPasswordToken(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	GetUserProfile(ctx *gin.Context)
	UpdateUserProfile(ctx *gin.Context)
	DeleteAccount(ctx *gin.Context)
}

// ServiceClient is a struct for the authentication service client
type ServiceClient struct {
	client pb_authentication.AuthenticationServiceClient
}

var _ ServiceClienter = &ServiceClient{}

// InitServiceClient initializes the authentication service client
func InitServiceClient(config *commonConfig.Config) (*ServiceClient, error) {
	grpcServiceAddress := fmt.Sprintf("%s:%s", config.AuthenticationService.Host, config.AuthenticationService.Port)

	fmt.Println("Connecting to authentication service at", grpcServiceAddress, config.TLSEnabled)
	clientConnection, err := commonTLS.CreateGRPCConnection(grpcServiceAddress, config.TLSEnabled)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to grpc authentication service: %v", err)
	}

	service := &ServiceClient{
		client: pb_authentication.NewAuthenticationServiceClient(clientConnection),
	}
	return service, nil
}

// GetPublicKey gets the public key from server
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

// Register redirects request to the register route
func (service *ServiceClient) Register(ctx *gin.Context) {
	authenticationroutes.Register(ctx, service.client)
}

// VerifyEmail redirects request to the verify email route
func (service *ServiceClient) VerifyEmail(ctx *gin.Context) {
	authenticationroutes.VerifyEmail(ctx, service.client)
}

// ResendEmailVerification redirects request to the resend email verification route
func (service *ServiceClient) ResendEmailVerification(ctx *gin.Context) {
	authenticationroutes.ResendEmailVerification(ctx, service.client)
}

// Authenticate redirects request to the authenticate route
func (service *ServiceClient) Authenticate(ctx *gin.Context) {
	authenticationroutes.Authenticate(ctx, service.client)
}

// AuthenticateWithFirebase redirects request to the authentication with firebase route
func (service *ServiceClient) AuthenticateWithFirebase(ctx *gin.Context) {
	authenticationroutes.AuthenticateWithFirebase(ctx, service.client)
}

// RefreshToken redirects request to the refresh token route
func (service *ServiceClient) RefreshToken(ctx *gin.Context) {
	authenticationroutes.RefreshToken(ctx, service.client)
}

// ForgotPassword redirects request to	the forgot password route
func (service *ServiceClient) ForgotPassword(ctx *gin.Context) {
	authenticationroutes.ForgotPassword(ctx, service.client)
}

// VerifyResetPasswordToken redirects request to the verify reset password token route
func (service *ServiceClient) VerifyResetPasswordToken(ctx *gin.Context) {
	authenticationroutes.VerifyResetPasswordToken(ctx, service.client)
}

// ResetPassword redirects request to the reset password route
func (service *ServiceClient) ResetPassword(ctx *gin.Context) {
	authenticationroutes.ResetPassword(ctx, service.client)
}

// GetUserProfile redirects request to the get user profile route
func (service *ServiceClient) GetUserProfile(ctx *gin.Context) {
	authenticationroutes.GetUserProfile(ctx, service.client)
}

// UpdateUserProfile redirects request to the update user profile route
func (service *ServiceClient) UpdateUserProfile(ctx *gin.Context) {
	authenticationroutes.UpdateUserProfile(ctx, service.client)
}

// DeleteAccount redirects request to the update user profile route
func (service *ServiceClient) DeleteAccount(ctx *gin.Context) {
	authenticationroutes.DeleteAccount(ctx, service.client)
}
