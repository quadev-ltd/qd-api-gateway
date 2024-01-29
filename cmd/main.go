package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	commonLogger "github.com/quadev-ltd/qd-common/pkg/log"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/authentication"
	"github.com/quadev-ltd/qd-qpi-gateway/internal/config"
)

func main() {
	configuration := config.Config{}
	err := configuration.Load("internal/config")
	if err != nil {
		log.Fatalln("Failed loading the configurations", err)
	}

	router := gin.Default()

	router.Use(commonLogger.AddNewCorrelationIDToContext)
	logger := commonLogger.NewLogFactory(configuration.Environment)
	router.Use(commonLogger.CreateGinLoggerMiddleware(logger))

	_, err = authentication.RegisterRoutes(router, &configuration)
	if err != nil {
		log.Fatalln("Failed to register authentication routes", err)
	}
	// // TODO: Add authentication middleware logic in routes which need to authenticate the user jwt
	// authenticationMiddleware, err := authentication.InitAuthenticationMiddleware(authenticationService)
	// if err != nil {
	// 	log.Fatalln("Failed to initiate authenticator middleware", err)
	// }
	// product.RegisterRoutes(router, &config, authSvc)
	// order.RegisterRoutes(router, &config, &authSvc)

	router.Run(fmt.Sprintf("%s:%s", configuration.GRPC.Host, configuration.GRPC.Port))
}
