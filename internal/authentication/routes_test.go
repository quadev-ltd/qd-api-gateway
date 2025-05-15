package authentication

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"
	"github.com/stretchr/testify/assert"

	authmock "github.com/quadev-ltd/qd-qpi-gateway/internal/authentication/mock"
)

func TestAuthenticateEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	router := gin.New()
	api := router.Group("/api/v1")

	// Create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock authentication service client
	mockClient := authmock.NewMockAuthenticationServiceClient(ctrl)

	// Create mock authentication middleware
	mockMiddleware := authmock.NewMockAutheticationMiddlewarer(ctrl)

	// Set up middleware expectations
	mockMiddleware.EXPECT().
		RequireAuthentication(gomock.Any()).
		DoAndReturn(func(c *gin.Context) {
			c.Next()
		}).
		AnyTimes()

	mockMiddleware.EXPECT().
		RefreshAuthentication(gomock.Any()).
		DoAndReturn(func(c *gin.Context) {
			c.Next()
		}).
		AnyTimes()

	// Create test service
	service := &ServiceClient{
		client: mockClient,
	}

	// Register routes
	err := RegisterRoutes(service, api, nil, nil, mockMiddleware)
	assert.NoError(t, err)

	// Test cases
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
		mockResponse   *pb_authentication.AuthenticateResponse
		mockError      error
	}{
		{
			name: "successful authentication",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"authToken":    "mock-token",
				"refreshToken": "mock-refresh-token",
			},
			mockResponse: &pb_authentication.AuthenticateResponse{
				AuthToken:    "mock-token",
				RefreshToken: "mock-refresh-token",
			},
			mockError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			mockClient.EXPECT().
				Authenticate(
					gomock.Any(),
					&pb_authentication.AuthenticateRequest{
						Email:    tt.requestBody["email"].(string),
						Password: tt.requestBody["password"].(string),
					},
				).
				Return(tt.mockResponse, tt.mockError)

			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/user/sessions", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}
