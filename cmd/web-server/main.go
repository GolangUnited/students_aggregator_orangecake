package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/indikator/aggregator_orange_cake/pkg/core"
)

func startServer(aLogger core.Logger, aServer *http.Server, aStopFunc context.CancelFunc) {
	defer aStopFunc()

	if aServer == nil {
		return
	}

	//TODO: should we print the following messages to console?
	// please note that we use console for the log output so the messages might be missed
	//fmt.Println("Starting server at " + aServer.Addr)
	//fmt.Println("Press Ctrl+C to exit")

	if lErr := aServer.ListenAndServe(); lErr != nil && !errors.Is(lErr, http.ErrServerClosed) {
		aLogger.Error("Cannot start http server. " + lErr.Error())
		return
	}

	aLogger.Info("Server has been stoped.")
}

func stopServer(aLogger core.Logger, aServer *http.Server) {
	aLogger.Info("Stopping http server")

	// Give some time to shutdown the server
	const SHUTDOWN_TIMEOUT = 30 * time.Second
	lShutdownCtx, lShutdownCancel := context.WithTimeout(context.Background(), SHUTDOWN_TIMEOUT)
	defer lShutdownCancel()

	if err := aServer.Shutdown(lShutdownCtx); err != nil {
		aLogger.Error("Shutdoun error." + err.Error())
	}
}

func main() {
	lFailed := false
	lLogger := initLogger(&lFailed)
	lConfig := initConfig(lLogger, &lFailed)
	lBuilder := initWebServerBuilder(lLogger, lConfig, &lFailed)
	lServer := initHttpServer(lBuilder, lConfig, &lFailed)
	if lFailed {
		return
	}

	lLogger.Info(fmt.Sprintf("Port: %d", lConfig.ServerPort))
	lLogger.Info(fmt.Sprintf("DB: %s", lConfig.DBConnectionString))

	lCtx, lStopFunc := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go startServer(lLogger, lServer, lStopFunc)

	<-lCtx.Done()
	lLogger.Trace("terminate signal received")
	stopServer(lLogger, lServer)

	lLogger.Info("Exit from aggregator's web server.")
}
