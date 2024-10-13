package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/config"
	internalhttp "github.com/hizagi/fullstack-challenge-jobsity/backend/internal/http"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/service"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage"
	"github.com/hizagi/fullstack-challenge-jobsity/backend/internal/storage/repository"
)

func main() {
	// root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failure loading configuration: %s", err)
	}

	httpServerConfig := serviceConfig.HTTPServerConfig()

	dbConfig, err := serviceConfig.DBConfig()
	if err != nil {
		log.Fatalf("failure loading database configuration: %s", err)
	}

	authConfig, err := serviceConfig.AuthConfig()
	if err != nil {
		log.Fatalf("failure loading auth configuration: %s", err)
	}

	mongoStorage, err := storage.NewMongoStorage(ctx, dbConfig)
	if err != nil {
		log.Fatalf("failure initializing mongodb storage: %s", err)
	}

	taskRepository := repository.NewTaskRepository(mongoStorage)

	taskService := service.NewTaskService(taskRepository)

	authMiddleware := internalhttp.APIKeyAuthMiddleware(authConfig.APIKey)

	httpHandler := internalhttp.NewTaskHandler(taskService, authMiddleware)

	srv := &http.Server{
		ReadTimeout:  httpServerConfig.ReadTimeout,
		WriteTimeout: httpServerConfig.WriteTimeout,
		IdleTimeout:  httpServerConfig.IdleTimeout,
		Addr:         fmt.Sprintf(":%d", httpServerConfig.Port),
		Handler:      httpHandler,
	}

	// Start the server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failure in the server: %s", err)
		}
	}()
	log.Printf("server started on :%d\n", httpServerConfig.Port)

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}

	log.Println("server exiting")
}
