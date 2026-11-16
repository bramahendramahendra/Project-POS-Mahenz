package main

import (
	"portal_api/config"
	error_helper "portal_api/helper/error"
	log_helper "portal_api/helper/log"
	"portal_api/pkg/database"
	"portal_api/routes"
	bootstrap "portal_api/server"

	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	engine := bootstrap.Initialized()
	port := fmt.Sprintf(":%v", config.ENV.AppPort)

	server := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	routes.PrintRoutes(engine, port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Cfg.MaxTimeoutGracefulShutdown)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		entry := log_helper.FromBackground("Internal Error", "Shutdown process", fmt.Sprintf("Shutdown Error: %s", err.Error()))
		entry.Stacktrace = error_helper.GetStackTrace(1)
		log_helper.LogError(entry)
	}

	if err := closeConnections(); err != nil {
		entry := log_helper.FromBackground("Internal Error", "Close all connections", fmt.Sprintf("Failed to close all connections: %s", err.Error()))
		entry.Stacktrace = error_helper.GetStackTrace(1)
		log_helper.LogError(entry)
	}

	log.Println("Server Exiting")
}

func closeConnections() error {
	var errs []string

	if err := database.DbManager.Close(config.Db.Database); err != nil {
		errs = append(errs, fmt.Sprintf("Failed to close database connection: %v", err))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, " | "))
	}

	log.Println("All connections closed successfully")
	return nil
}
