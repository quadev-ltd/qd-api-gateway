package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	commonPB "github.com/quadev-ltd/qd-common/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCErrorToHTTPStatus converts a gRPC error to an HTTP status code
func GRPCErrorToHTTPStatus(err error) int {
	st, ok := status.FromError(err)
	if !ok {
		// If the error is not a gRPC status error, default to 500
		return http.StatusInternalServerError
	}

	switch st.Code() {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition, codes.Aborted:
		return http.StatusPreconditionFailed
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// HandleError handles an error by returning an HTTP response with the appropriate status code
func HandleError(ctx *gin.Context, err error) error {
	errorHTTPStatusCode := GRPCErrorToHTTPStatus(err)
	errorsMap := gin.H{"error": err.Error()}

	fieldValidationErrors, parsingError := commonPB.GetFieldValidationErrors(err)
	if parsingError != nil {
		ctx.JSON(errorHTTPStatusCode, errorsMap)
		ctx.AbortWithError(errorHTTPStatusCode, err)
		return parsingError
	}

	if len(fieldValidationErrors) > 0 {
		errorsMap["field_errors"] = fieldValidationErrors
	}
	ctx.JSON(errorHTTPStatusCode, errorsMap)
	ctx.AbortWithError(errorHTTPStatusCode, err)
	return nil
}
