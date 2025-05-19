package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	commontConfig "github.com/quadev-ltd/qd-common/pkg/config"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/services"
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

	serviceInitializer := services.NewServiceInitialiser(&configuration, &centralConfig, router, api)
	if err := serviceInitializer.InitializeAllServices(); err != nil {
		log.Fatalln("Failed to initialize services:", err)
	}

	fmt.Println("Listening API requests on URL: ", fmt.Sprintf("%s:%s%s", centralConfig.GatewayService.Host, centralConfig.GatewayService.Port, APIPath))
	router.Run(fmt.Sprintf("%s:%s", centralConfig.GatewayService.Host, centralConfig.GatewayService.Port))
}
