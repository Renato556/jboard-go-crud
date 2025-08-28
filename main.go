package main

import (
	"context"
	httpt "jboard-go-crud/src/controllers"
	"jboard-go-crud/src/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"jboard-go-crud/src/config"
	mongoadapter "jboard-go-crud/src/repositories"
	"jboard-go-crud/src/services"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// 1) Infra: Mongo client
	config.InitConnection()
	client := config.GetMongoClient()

	dbName := os.Getenv("MONGODB_DATABASE_NAME")
	collName := os.Getenv("MONGODB_JOB_COLLECTION")

	// 2) Repositório
	jobRepo := mongoadapter.NewJobRepository(client, dbName, collName)

	// 3) Serviço
	jobService := services.NewJobService(jobRepo)

	// 4) HTTP Handler + Router
	jobHandler := httpt.NewJobHandler(jobService)
	router := routers.NewJobsController(jobHandler)

	// 5) Servidor HTTP
	srv := &http.Server{
		Addr:              ":8081",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Server started at port 8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("erro ao iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
	_ = config.CloseConnection(ctx)
}
