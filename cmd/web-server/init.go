package main

import (
	"fmt"
	"net/http"
	"os"

	http_handlers "github.com/indikator/aggregator_orange_cake/cmd/web-server/handlers"
	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

func initLogger(aFailed *bool) core.Logger {
	if *aFailed {
		return nil
	}

	lLogger := core.NewZeroLogger(os.Stdout)
	lLogger.Info("Starting the aggregator's web server.")

	return lLogger
}

func initConfig(aLogger core.Logger, aFailed *bool) *ServiceConfig {
	if *aFailed {
		return nil
	}

	lConfig, lErr := NewServiceConfig()
	if lErr != nil {
		*aFailed = true
		aLogger.Error(fmt.Sprintf("Config Error: %s\n", lErr.Error()))
		return nil
	}

	return lConfig
}

func initWebServerBuilder(aLogger core.Logger, aConfig *ServiceConfig, aFailed *bool) http_handlers.WebServerBuilder {
	if *aFailed {
		return nil
	}

	//TODO: use an SqliteStorage storage
	lStorage := core.MockDBStorage{}
	/*
		if lErr != nil {
			*aFailed = true
			aLogger.Error(fmt.Sprintf("Init storage Error: %s\n", lErr.Error()))
			return nil
		}
	*/

	lWebServer, lErr := http_handlers.NewWebServer(aLogger, lStorage, lStorage)
	if lErr != nil {
		*aFailed = true
		aLogger.Error(fmt.Sprintf("Init web server Error: %s\n", lErr.Error()))
		return nil
	}

	return http_handlers.NewWebServerBuilder(lWebServer)
}

func initHttpServer(aBuilder http_handlers.WebServerBuilder, aConfig *ServiceConfig, aFailed *bool) *http.Server {
	if *aFailed {
		return nil
	}

	aBuilder.OpenApi().Spec.Info.
		WithTitle("Test Aggregator API").
		WithVersion("1.0.0").
		WithDescription("Will write some description later")

	http_handlers.RegisterGetArticles("/api/articles", aBuilder, aFailed)

	// open API handlers must be the last ones (all other handlers must be registered in OpenApi already)
	http_handlers.RegisterOpenApiJson("/openapi.json", aBuilder, aFailed)
	http_handlers.RegisterOpenApiUI("/", aBuilder, aFailed)

	if *aFailed {
		return nil
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", aConfig.ServerPort),
		Handler: aBuilder.ServeMux(),
	}
}
