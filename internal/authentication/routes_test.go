package authentication

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb_authentication.RegisterAuthenticationServiceServer(s, &mockAuthServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("Server exited with error: %v", err))
		}
	}()
}

// mockAuthServer implements the AuthenticationServiceServer interface
type mockAuthServer struct {
	pb_authentication.UnimplementedAuthenticationServiceServer
	lastAuthenticateRequest *pb_authentication.AuthenticateRequest
}

func (s *mockAuthServer) Authenticate(ctx context.Context, req *pb_authentication.AuthenticateRequest) (*pb_authentication.AuthenticateResponse, error) {
	s.lastAuthenticateRequest = req
	return &pb_authentication.AuthenticateResponse{
		AuthToken:    "mock-auth-token",
		RefreshToken: "mock-refresh-token",
	}, nil
}

func (s *mockAuthServer) GetPublicKey(ctx context.Context, req *pb_authentication.GetPublicKeyRequest) (*pb_authentication.GetPublicKeyResponse, error) {
	return &pb_authentication.GetPublicKeyResponse{
		PublicKey: "mock-public-key",
	}, nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestAuthenticateEndpointIntegration(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(commonLogger.AddNewCorrelationIDToContext)

	// Create mock server
	mockServer := &mockAuthServer{}

	// Set up routes
	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/sessions", func(c *gin.Context) {
				var req struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// Call mock server
				resp, err := mockServer.Authenticate(c.Request.Context(), &pb_authentication.AuthenticateRequest{
					Email:    req.Email,
					Password: req.Password,
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"auth_token":    resp.AuthToken,
					"refresh_token": resp.RefreshToken,
				})
			})
		}
	}

	// Create test server
	ts := &testServer{
		router:         router,
		mockAuthServer: mockServer,
	}

	// Start server
	go ts.start()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test cases
	t.Run("successful_authentication", func(t *testing.T) {
		// Make HTTP request to authenticate endpoint
		resp, err := ts.makeRequest("POST", "/api/v1/user/sessions", map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		})
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		// Verify gRPC request parameters
		assert.NotNil(t, ts.mockAuthServer.lastAuthenticateRequest)
		assert.Equal(t, "test@example.com", ts.mockAuthServer.lastAuthenticateRequest.Email)
		assert.Equal(t, "password123", ts.mockAuthServer.lastAuthenticateRequest.Password)
	})
}

// testServer represents a test HTTP server
type testServer struct {
	router         *gin.Engine
	mockAuthServer *mockAuthServer
}

func (ts *testServer) start() {
	ts.router.Run(":8080")
}

func (ts *testServer) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:8080%s", path), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
