package routers

import (
	"jboard-go-crud/internal/controllers"
	"net/http"
)

func NewJobsController(jobHandler *controllers.JobHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/jobs", jobHandler.CreateJob)
	mux.HandleFunc("GET /v1/jobs", jobHandler.GetAllJobs)
	mux.HandleFunc("GET /v1/health", healthCheck)
	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy","service":"jboard-go-crud"}`))
}
