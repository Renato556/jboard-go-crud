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
		log.Printf("Invalid JSON payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Log apenas campos críticos para debugging
	log.Printf("Job payload: ID='%s', Title='%s', Field='%s'", job.ID, job.Title, job.Field)

	outcome, err := h.svc.CreateOrUpdate(r.Context(), job)
	if err != nil {
		log.Printf("CreateOrUpdate failed for job '%s': %v", job.ID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch outcome {
	case services.OutcomeCreated:
		log.Printf("Job created: %s", job.ID)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "Job created successfully.",
		})
	case services.OutcomeUpdated:
		log.Printf("Job updated: %s", job.ID)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "Job already exists, updated with new information and extended expiration.",
		})
	default:
		log.Printf("Unknown outcome %d for job %s", outcome, job.ID)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
	}
}

func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "HTTP Method invalid", http.StatusMethodNotAllowed)
		return
	}

	jobs, err := h.svc.FindAll(r.Context())
	if err != nil {
		log.Printf("FindAll failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		log.Printf("JSON encode error: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("Returned %d jobs", len(jobs))
}
