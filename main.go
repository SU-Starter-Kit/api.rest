package main

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/su-starter-kit/api.rest/app/config"
	v1 "github.com/su-starter-kit/api.rest/app/v1"
	"github.com/su-starter-kit/api.rest/docs"
	"github.com/su-starter-kit/log/logger"
	"github.com/su-starter-kit/log/messages"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	lggr := buildLogger("None")
	cfg := config.NewAppConfig().WithEnvVarGetter()
	port := fmt.Sprintf(":%s", cfg.GetParameter("PORT"))

	// Sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://1f4807f97cdd40059b73be6f24e6d25c@o4504168838070272.ingest.sentry.io/4504170441146368",
		TracesSampleRate: 1.0,
	}); err != nil {
		lggr.Error(&messages.M{Message: "Error while starting sentry", Err: err.Error()})
	}

	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")

	router := buildRouter(lggr)
	{
		lggr.Info(&messages.M{Message: fmt.Sprintf("Starting at port: %s", port)})
		sentry.CaptureException(fmt.Errorf("Error!"))

		if err := router.Run(port); err != nil {
			lggr.Warn(&messages.M{Message: "Error while starting router"})
		}
	}
}

func buildRouter(lggr logger.CompanyLogger) *gin.Engine {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/v1/user"

	v1UserHandler := v1.NewUserHandler(lggr)

	// v1User routes
	v1User := router.Group("/v1/user")
	{
		v1User.POST("/login", v1UserHandler.UserLogin)
		v1User.POST("/submit", nil)
	}

	// Swagger routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}

func buildLogger(correlationId string) logger.CompanyLogger {
	return logger.NewBuilder().
		WithCorrelationId(correlationId).
		Build()
}
