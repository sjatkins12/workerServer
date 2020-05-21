package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sjatkins12/workerServer/router"
	"github.com/sjatkins12/workerServer/utils"
	"github.com/sjatkins12/workerServer/worker"
)

func main() {
	utils.Log.Info("Starting microservice server")

	taskQueue := make(chan worker.Job, 10)
	shutdown := make(chan struct{})

	go worker.Worker(taskQueue, shutdown)

	router := router.SetupRouter(taskQueue)

	srv := &http.Server{
		Addr:    ":8017",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Log.Error("Error while listening to http server")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until TERM signal
	<-quit
	close(shutdown)

	utils.Log.Info("Shutdown pipeline ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		utils.Log.Error("Server Shutdown Error: ", err)
		os.Exit(1)
	}

	utils.Log.Info("server shutdown gracefully")
}
