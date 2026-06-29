package main

import (
	"pos_api/config"
	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
	cash_drawer_service "pos_api/domain/cash_drawer/service"
	error_helper "pos_api/helper/error"
	log_helper "pos_api/helper/log"
	"pos_api/pkg/database"
	pkgdatabase "pos_api/pkg/database"
	"pos_api/routes"
	"pos_api/scheduler"
	bootstrap "pos_api/server"

	// "pos_api/pkg/redis"

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

	schedulerCtx, cancelScheduler := context.WithCancel(context.Background())
	cdRepo := cash_drawer_repo.NewCashDrawerRepo(pkgdatabase.DB)
	cdService := cash_drawer_service.NewCashDrawerService(cdRepo)
	scheduler.StartCashDrawerScheduler(schedulerCtx, cdService)

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
	cancelScheduler()

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

	// If you enable Redis later
	// if err := redis.Client.Close(); err != nil {
	//     errs = append(errs, fmt.Sprintf("Failed to close redis connection: %v", err))
	// }

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, " | "))
	}

	log.Println("All connections closed successfully")
	return nil
}
