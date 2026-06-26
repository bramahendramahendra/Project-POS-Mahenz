package main

import (
	"pos_api/config"
	cash_drawer_repo "pos_api/domain/cash_drawer/repo"
	cash_drawer_service "pos_api/domain/cash_drawer/service"
	"pos_api/helper"
	error_helper "pos_api/helper/error"
	time_helper "pos_api/helper/time"
	"pos_api/pkg/database"
	pkgdatabase "pos_api/pkg/database"
	"pos_api/pkg/logger"
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

	startTime := time_helper.GetTimeWithFormat()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Cfg.MaxTimeoutGracefulShutdown)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Error(fmt.Sprintf("Shutdown Error: %s", err.Error()),
			"Internal Error", "Shutdown process", helper.GenerateUniqueId(),
			error_helper.GetStackTrace(1), startTime, time_helper.GetEndTime(startTime), nil)
	}

	if err := closeConnections(); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to close all connections: %s", err.Error()),
			"Internal Error", "Close all connections", helper.GenerateUniqueId(),
			error_helper.GetStackTrace(1), startTime, time_helper.GetEndTime(startTime), nil)
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
