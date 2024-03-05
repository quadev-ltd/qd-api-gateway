package authentication

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"
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

func transferCorrelationIDToOutgoingContext(ctx *gin.Context) {
	logger, err := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	contextWithCorrelationID, err := commonLogger.TransferCorrelationIDToOutgoingContext(ctx.Request.Context())
	if err != nil {
		logger.Error(err, "Error transferring correlation ID to outgoing context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	logger.Info("Correlation ID successfully transfered to outgoing context")
	ctx.Request = ctx.Request.WithContext(contextWithCorrelationID)
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
	transferCorrelationIDToOutgoingContext(ctx)
	routes.Register(ctx, service.client)
}

func (service *ServiceClient) VerifyEmail(ctx *gin.Context) {
	transferCorrelationIDToOutgoingContext(ctx)
	routes.VerifyEmail(ctx, service.client)
}

func (service *ServiceClient) ResendEmailVerification(ctx *gin.Context) {
	transferCorrelationIDToOutgoingContext(ctx)
	routes.ResendEmailVerification(ctx, service.client)
}

func (service *ServiceClient) Authenticate(ctx *gin.Context) {
	transferCorrelationIDToOutgoingContext(ctx)
	routes.Authenticate(ctx, service.client)
}

func (service *ServiceClient) RefreshToken(ctx *gin.Context) {
	transferCorrelationIDToOutgoingContext(ctx)
	routes.RefreshToken(ctx, service.client)
}
