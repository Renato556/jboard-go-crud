package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"jboard-go-crud/internal/models"
	"jboard-go-crud/internal/services"
)

type JobHandler struct {
	svc services.JobService
}

func NewJobHandler(s services.JobService) *JobHandler {
	return &JobHandler{svc: s}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "HTTP Method invalid", http.StatusMethodNotAllowed)
		return
	}

	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		log.Println(err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	log.Printf("POST request received at /jobs with payload: %v", job)

	outcome, err := h.svc.CreateOrUpdate(r.Context(), job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch outcome {
	case services.OutcomeCreated:
		log.Printf("Job created successfully")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "Job created successfully.",
		})
	case services.OutcomeUpdated:
		log.Printf("Job already exists, updated with new information and extended expiration")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "Job already exists, updated with new information and extended expiration.",
		})
	default:
		http.Error(w, "Unknown error", http.StatusInternalServerError)
	}
}

func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "HTTP Method invalid", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("GET request received at /v1/jobs")

	jobs, err := h.svc.FindAll(r.Context())
	if err != nil {
		log.Printf("Error fetching jobs: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully returned %d jobs", len(jobs))
}
