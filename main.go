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

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	log.Printf("Starting jboard-go-crud application...")

	_ = godotenv.Load()

	// 1) MongoDB connection
	config.InitConnection()
	client := config.GetClient()

	dbName := os.Getenv("MONGODB_DATABASE_NAME")
	if dbName == "" {
		log.Fatal("MONGODB_DATABASE_NAME environment variable not set")
	}

	// 2) Initialize collections configuration
	jobCollName := os.Getenv("MONGODB_JOB_COLLECTION")
	userCollName := os.Getenv("MONGODB_USER_COLLECTION")

	// Set default collection names if not provided
	if jobCollName == "" {
		jobCollName = "jobs"
		log.Printf("Using default jobs collection name: %s", jobCollName)
	}
	if userCollName == "" {
		userCollName = "users"
		log.Printf("Using default users collection name: %s", userCollName)
	}

	log.Printf("MongoDB configuration - DB: '%s', Jobs Collection: '%s', Users Collection: '%s'",
		dbName, jobCollName, userCollName)

	// 3) Initialize repositories and services
	jobRepo := repositories.NewJobRepository(client, dbName, jobCollName)
	jobService := services.NewJobService(jobRepo)
	jobHandler := controllers.NewJobHandler(jobService)

	userRepo := repositories.NewUserRepository(client, dbName, userCollName)
	userService := services.NewUserService(userRepo)
	userHandler := controllers.NewUserHandler(userService)

	skillRepo := repositories.NewSkillRepository(client, dbName, "skills")
	skillService := services.NewSkillService(skillRepo)
	skillHandler := controllers.NewSkillHandler(skillService)

	// 4) Initialize routers
	jobRouter := routers.NewJobsController(jobHandler)
	userRouter := routers.NewUsersController(userHandler)
	skillRouter := routers.NewSkillsController(skillHandler)

	// 5) Create main router and mount sub-routers
	mainRouter := mux.NewRouter()
	mainRouter.PathPrefix("/v1/jobs").Handler(jobRouter)
	mainRouter.PathPrefix("/v1/users").Handler(userRouter)
	mainRouter.PathPrefix("/v1/skills").Handler(skillRouter)

	// 6) HTTP Server
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mainRouter,
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
