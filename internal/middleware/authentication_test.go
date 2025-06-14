package middleware

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	commmonJWT "github.com/quadev-ltd/qd-common/pkg/jwt"
	commonJWT "github.com/quadev-ltd/qd-common/pkg/jwt"
	commonJWTMock "github.com/quadev-ltd/qd-common/pkg/jwt/mock"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"
	commonLoggerMock "github.com/quadev-ltd/qd-common/pkg/log/mock"
	commonToken "github.com/quadev-ltd/qd-common/pkg/token"
	"github.com/stretchr/testify/assert"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/authentication/mock"
)

func createTestContext(method, path string, body []byte, authHeader *string) (*gin.Context, *httptest.ResponseRecorder) {
	// Initialize Gin engine
	gin.SetMode(gin.TestMode)

	// Create a request
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// You can add additional setup here, such as setting headers
	c.Request.Header.Set("Content-Type", "application/json")
	if authHeader != nil {
		c.Request.Header.Set("Authorization", *authHeader)
	}

	return c, w
}

func createTestContextWithLogger(logger commonLogger.Loggerer, authHeader *string) (*gin.Context, *httptest.ResponseRecorder) {
	ctx, w := createTestContext("GET", "/test", nil, authHeader)
	newCtx := context.WithValue(ctx.Request.Context(), commonLogger.LoggerKey, logger)
	newReq := ctx.Request.WithContext(newCtx)
	ctx.Request = newReq
	return ctx, w
}

func fastBackoff(attempt int) time.Duration {
	return 10 * time.Millisecond // or time.Duration(0) for no delay
}

func TestMiddleware(t *testing.T) {
	environment := "Test"

	// RequestPublicKey
	t.Run("Request_Public_Key_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		errorExample := errors.New("example error")
		correlationID := "example-correlation-id"

		serviceMock.EXPECT().GetPublicKey(gomock.Any()).Return(nil, errorExample).Times(5)

		publicKey, err := RequestPublicKey(serviceMock, correlationID, environment, fastBackoff)

		assert.Error(t, err)
		assert.Nil(t, publicKey)
		assert.Equal(t, "Could not obtain public key after 5 attempts: example error", err.Error())
	})

	t.Run("Request_Public_Key_Success", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		correlationID := "example-correlation-id"
		publicKeyExample := "example-key"

		serviceMock.EXPECT().GetPublicKey(gomock.Any()).Return(&publicKeyExample, nil)

		publicKey, err := RequestPublicKey(serviceMock, correlationID, environment, fastBackoff)

		assert.Nil(t, err)
		assert.Equal(t, *publicKey, publicKeyExample)
	})

	// RequireAuthentication
	t.Run("RequireAuthentication_No_Logger_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		ctx, w := createTestContext("GET", "/test", nil, nil)

		authenticationMiddleware.RequireAuthentication(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("RequireAuthentication_No_Authorization_Header_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		ctx, w := createTestContextWithLogger(loggerMock, nil)

		loggerMock.EXPECT().Error(nil, "No authorization header was present in the request")

		authenticationMiddleware.RequireAuthentication(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("RequireAuthentication_Wrong_Authorization_Header_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		authHeader := "test-header"
		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		loggerMock.EXPECT().Error(nil, "No bearer token was present in the authorization header")

		authenticationMiddleware.RequireAuthentication(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("RequireAuthentication_Empty_Authorization_Header_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		authHeader := "Bearer"
		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		loggerMock.EXPECT().Error(nil, "No bearer token was present in the authorization header")

		authenticationMiddleware.RequireAuthentication(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("RequireAuthentication_Invalid_Authorization_Header_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		exampleError := errors.New("example error")
		authHeader := "Bearer invalid-header"
		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		jwtVerifierMock.EXPECT().Verify("invalid-header").Return(nil, exampleError)

		authenticationMiddleware.RequireAuthentication(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("RequireAuthentication_Type_Claim_Authorization_Header_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		exampleError := errors.New("example error")
		testTokenValue := "test-header"
		authHeader := "Bearer " + testTokenValue
		testToken := &jwt.Token{}

		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		jwtVerifierMock.EXPECT().Verify(testTokenValue).Return(testToken, nil)
		jwtTokenInspectorMock.EXPECT().GetClaimsFromToken(testToken).Return(nil, exampleError)

		authenticationMiddleware.RequireAuthentication(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("RequireAuthentication_Wrong_Type_Claim_Authorization_Header_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		authHeader := "Bearer test-header"
		testToken := &jwt.Token{}
		tokenClaims := &commmonJWT.TokenClaims{
			Type: commonToken.RefreshTokenType,
		}

		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		jwtVerifierMock.EXPECT().Verify("test-header").Return(testToken, nil)
		jwtTokenInspectorMock.EXPECT().GetClaimsFromToken(testToken).Return(tokenClaims, nil)

		authenticationMiddleware.RequireAuthentication(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("RequireAuthentication_Expiry_Claim_Authorization_Header_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		authHeader := "Bearer test-header"
		testToken := &jwt.Token{}
		tokenClaims := &commmonJWT.TokenClaims{
			Type:   commonToken.AuthTokenType,
			Expiry: time.Now().Add(-1 * time.Second),
		}

		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		jwtVerifierMock.EXPECT().Verify("test-header").Return(testToken, nil)
		jwtTokenInspectorMock.EXPECT().GetClaimsFromToken(testToken).Return(tokenClaims, nil)

		authenticationMiddleware.RequireAuthentication(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("RequireAuthentication_Success", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		authHeader := "Bearer test-header"
		testToken := &jwt.Token{}
		tokenClaims := &commmonJWT.TokenClaims{
			Type:   commonToken.AuthTokenType,
			Expiry: time.Now().Add(10 * time.Second),
		}

		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		jwtVerifierMock.EXPECT().Verify("test-header").Return(testToken, nil)
		jwtTokenInspectorMock.EXPECT().GetClaimsFromToken(testToken).Return(tokenClaims, nil)
		loggerMock.EXPECT().Info("Successfully authenticated user")

		authenticationMiddleware.RequireAuthentication(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Refresh Authentication
	t.Run("RefreshAuthentication_Wrong_Type_Claim_Authorization_Header_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		authHeader := "Bearer test-header"
		testToken := &jwt.Token{}
		tokenClaims := &commmonJWT.TokenClaims{
			Type:   commonToken.AuthTokenType,
			Expiry: time.Now().Add(-1 * time.Second),
		}

		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		jwtVerifierMock.EXPECT().Verify("test-header").Return(testToken, nil)
		jwtTokenInspectorMock.EXPECT().GetClaimsFromToken(testToken).Return(tokenClaims, nil)

		authenticationMiddleware.RefreshAuthentication(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("RequirePaidFeatures_Success", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		authHeader := "Bearer test-header"
		testToken := &jwt.Token{}
		tokenClaims := &commmonJWT.TokenClaims{
			Type:            commonToken.AuthTokenType,
			Expiry:          time.Now().Add(10 * time.Second),
			HasPaidFeatures: true,
		}
		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		jwtVerifierMock.EXPECT().Verify("test-header").Return(testToken, nil)
		jwtTokenInspectorMock.EXPECT().GetClaimsFromToken(testToken).Return(tokenClaims, nil)
		jwtTokenInspectorMock.EXPECT().GetClaimsFromToken(testToken).Return(tokenClaims, nil)
		loggerMock.EXPECT().Info("Successfully authenticated user")
		ctx.Set(string(commonJWT.JWTTokenKey), testToken)
		loggerMock.EXPECT().Info("User has access to paid features")

		authenticationMiddleware.RequirePaidFeatures(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("RequirePaidFeatures_No_Paid_Features_Error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		serviceMock := mock.NewMockServiceClienter(controller)
		jwtVerifierMock := commonJWTMock.NewMockTokenVerifierer(controller)
		jwtTokenInspectorMock := commonJWTMock.NewMockTokenInspectorer(controller)
		authenticationMiddleware := &AutheticationMiddleware{
			serviceMock,
			jwtVerifierMock,
			jwtTokenInspectorMock,
		}
		loggerMock := commonLoggerMock.NewMockLoggerer(controller)

		authHeader := "Bearer test-header"
		testToken := &jwt.Token{}
		tokenClaims := &commmonJWT.TokenClaims{
			Type:            commonToken.AuthTokenType,
			Expiry:          time.Now().Add(10 * time.Second),
			HasPaidFeatures: false,
		}

		ctx, w := createTestContextWithLogger(loggerMock, &authHeader)

		jwtVerifierMock.EXPECT().Verify("test-header").Return(testToken, nil)
		jwtTokenInspectorMock.EXPECT().GetClaimsFromToken(testToken).Return(tokenClaims, nil)
		jwtTokenInspectorMock.EXPECT().GetClaimsFromToken(testToken).Return(tokenClaims, nil)
		loggerMock.EXPECT().Info("Successfully authenticated user")
		ctx.Set(string(commonJWT.JWTTokenKey), testToken)

		authenticationMiddleware.RequirePaidFeatures(ctx)

		assert.Equal(t, http.StatusPaymentRequired, w.Code)
	})
}
