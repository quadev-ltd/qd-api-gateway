package authentication

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	commonConfig "github.com/quadev-ltd/qd-common/pkg/config"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"
	commonTLS "github.com/quadev-ltd/qd-common/pkg/tls"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/authentication/routes"
	"github.com/quadev-ltd/qd-qpi-gateway/pb/gen/go/pb_authentication"
)

type ServiceClient struct {
	Client pb_authentication.AuthenticationServiceClient
}

func InitServiceClient(config *commonConfig.Config) (pb_authentication.AuthenticationServiceClient, error) {
	grpcServiceAddress := fmt.Sprintf("%s:%s", config.AuthenticationService.Host, config.AuthenticationService.Port)

	clientConnection, err := commonTLS.CreateGRPCConnection(grpcServiceAddress, config.TLSEnabled)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to grpc authentication service: %v", err)
	}

	return pb_authentication.NewAuthenticationServiceClient(clientConnection), nil
}

func (service *ServiceClient) Register(ctx *gin.Context) {
	logger := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	contextWithCorrelationID, err := commonLogger.TransferCorrelationIDToOutgoingContext(ctx.Request.Context())
	if err != nil {
		logger.Error(err, "Error transferring correlation ID to outgoing context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Request = ctx.Request.WithContext(contextWithCorrelationID)
	routes.Register(ctx, service.Client)
}

func (service *ServiceClient) VerifyEmail(ctx *gin.Context) {
	logger := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	contextWithCorrelationID, err := commonLogger.TransferCorrelationIDToOutgoingContext(ctx.Request.Context())
	if err != nil {
		logger.Error(err, "Error transferring correlation ID to outgoing context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Request = ctx.Request.WithContext(contextWithCorrelationID)
	routes.VerifyEmail(ctx, service.Client)
}

func (service *ServiceClient) ResendEmailVerification(ctx *gin.Context) {
	logger := commonLogger.GetLoggerFromContext(ctx.Request.Context())
	contextWithCorrelationID, err := commonLogger.TransferCorrelationIDToOutgoingContext(ctx.Request.Context())
	if err != nil {
		logger.Error(err, "Error transferring correlation ID to outgoing context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Request = ctx.Request.WithContext(contextWithCorrelationID)
	routes.ResendEmailVerification(ctx, service.Client)
}
