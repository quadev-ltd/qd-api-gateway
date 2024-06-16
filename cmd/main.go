package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	commontConfig "github.com/quadev-ltd/qd-common/pkg/config"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/authentication"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
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

	// Universal links
	router.StaticFile("/.well-known/apple-app-site-association", "./.well-known/apple-app-site-association")
	router.StaticFile("/.well-known/assetlinks.json", "./.well-known/assetlinks.json")

	api := router.Group(APIPath)

	_, err = authentication.RegisterRoutes(api, &centralConfig, &configuration)
	if err != nil {
		log.Fatalln("Failed to register authentication routes: ", err)
	}
	fmt.Println("Listening API requests on URL: ", fmt.Sprintf("%s:%s%s", centralConfig.GatewayService.Host, centralConfig.GatewayService.Port, APIPath))
	router.Run(fmt.Sprintf("%s:%s", centralConfig.GatewayService.Host, centralConfig.GatewayService.Port))
}
