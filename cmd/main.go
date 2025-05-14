package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	commontConfig "github.com/quadev-ltd/qd-common/pkg/config"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/authentication"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/image_analysis"
	sharedMiddleware "github.com/quadev-ltd/qd-qpi-gateway/internal/shared/middleware"
)

// APIPath is the path of the API
const APIPath = "/api/v1"

func main() {
	configuration := config.Config{}
	err := configuration.Load("internal/config")
	if err != nil {
		log.Fatalln("Failed loading the configurations", err)
	}

	var centralConfig commontConfig.Config
	centralConfig.Load(
		configuration.Environment,
		configuration.AWS.Key,
		configuration.AWS.Secret,
	)

	router := gin.Default()
	router.Use(commonLogger.AddNewCorrelationIDToContext)
	logger := commonLogger.NewLogFactory(configuration.Environment)
	router.Use(commonLogger.CreateGinLoggerMiddleware(logger))

	api := router.Group(APIPath)

	authService, err := authentication.RegisterRoutes(api, &centralConfig, &configuration)
	if err != nil {
		log.Fatalln("Failed to register authentication routes: ", err)
	}

	authMiddleware, err := authentication.InitAuthenticationMiddleware(authService, &configuration)
	if err != nil {
		log.Fatalln("Failed to initialize authentication middleware: ", err)
	}

	_, err = image_analysis.RegisterRoutes(api, &configuration, authMiddleware.(sharedMiddleware.AutheticationMiddlewarer))
	if err != nil {
		log.Fatalln("Failed to register image analysis routes: ", err)
	}

	fmt.Println("Listening API requests on URL: ", fmt.Sprintf("%s:%s%s", centralConfig.GatewayService.Host, centralConfig.GatewayService.Port, APIPath))
	router.Run(fmt.Sprintf("%s:%s", centralConfig.GatewayService.Host, centralConfig.GatewayService.Port))
}
