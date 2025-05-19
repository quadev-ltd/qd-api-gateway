package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	commontConfig "github.com/quadev-ltd/qd-common/pkg/config"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/authentication"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/imageanalysis"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/middleware"
)

// ServiceInitialiser handles initialization of all services
type ServiceInitialiser struct {
	config               *config.Config
	centralConfig        *commontConfig.Config
	router               *gin.Engine
	apiGroup             *gin.RouterGroup
	authMiddleware       middleware.AutheticationMiddlewarer
	authService          authentication.ServiceClienter
	imageAnalysisService imageanalysis.ServiceClienter
}

// NewServiceInitialiser creates a new ServiceInitializer
func NewServiceInitialiser(config *config.Config, centralConfig *commontConfig.Config, router *gin.Engine, apiGroup *gin.RouterGroup) *ServiceInitialiser {
	return &ServiceInitialiser{
		config:        config,
		centralConfig: centralConfig,
		router:        router,
		apiGroup:      apiGroup,
	}
}

// InitializeAuthService initializes the authentication service and middleware
func (serviceInitialiser *ServiceInitialiser) InitializeAuthService() error {
	authService, err := authentication.InitServiceClient(serviceInitialiser.centralConfig)
	if err != nil {
		return fmt.Errorf("could not initialize authentication service client: %w", err)
	}
	serviceInitialiser.authService = authService

	authMiddleware, err := middleware.InitAuthenticationMiddleware(authService, serviceInitialiser.config)
	if err != nil {
		return fmt.Errorf("failed to initiate authenticator middleware: %w", err)
	}
	serviceInitialiser.authMiddleware = authMiddleware

	err = authentication.RegisterRoutes(authService, serviceInitialiser.apiGroup, serviceInitialiser.centralConfig, serviceInitialiser.config, authMiddleware)
	if err != nil {
		return fmt.Errorf("failed to register authentication routes: %w", err)
	}

	return nil
}

// InitializeImageAnalysisService initializes the image analysis service
func (serviceInitialiser *ServiceInitialiser) InitializeImageAnalysisService() error {
	if serviceInitialiser.authMiddleware == nil {
		return fmt.Errorf("authentication middleware not initialized")
	}

	imageAnalysisService, err := imageanalysis.InitServiceClient(serviceInitialiser.centralConfig)
	if err != nil {
		return fmt.Errorf("could not initialize image analysis service client: %w", err)
	}
	serviceInitialiser.imageAnalysisService = imageAnalysisService

	err = imageanalysis.RegisterRoutes(imageAnalysisService, serviceInitialiser.apiGroup, serviceInitialiser.centralConfig, serviceInitialiser.authMiddleware)
	if err != nil {
		return fmt.Errorf("failed to register image analysis routes: %w", err)
	}

	return nil
}

// InitializeAllServices initializes all services
func (serviceInitialiser *ServiceInitialiser) InitializeAllServices() error {
	if err := serviceInitialiser.InitializeAuthService(); err != nil {
		return err
	}

	if err := serviceInitialiser.InitializeImageAnalysisService(); err != nil {
		return err
	}

	return nil
}
