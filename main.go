package main

import (
	"context"
	"jboard-go-crud/internal/config"
	"jboard-go-crud/internal/controllers"
	"jboard-go-crud/internal/repositories"
	"jboard-go-crud/internal/routers"
	"jboard-go-crud/internal/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	log.Printf("Starting jboard-go-crud application...")

	_ = godotenv.Load()

	// 1) MongoDB connection
	config.InitConnection()
	client := config.GetClient()

	dbName := os.Getenv("MONGODB_DATABASE_NAME")
	collName := os.Getenv("MONGODB_JOB_COLLECTION")

	if dbName == "" || collName == "" {
		log.Printf("WARNING: Missing MongoDB environment variables - DB: '%s', Collection: '%s'", dbName, collName)
	}

	// 2) Initialize layers
	jobRepo := repositories.NewJobRepository(client, dbName, collName)
	jobService := services.NewJobService(jobRepo)
	jobHandler := controllers.NewJobHandler(jobService)
	router := routers.NewJobsController(jobHandler)

	// 3) HTTP Server
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Server ready at port 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Printf("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	config.CloseConnection()

	log.Printf("Application stopped")
}
