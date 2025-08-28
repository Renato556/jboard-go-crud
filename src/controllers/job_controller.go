package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"jboard-go-crud/src/models"
	"jboard-go-crud/src/services"
)

type JobHandler struct {
	svc services.JobService
}

func NewJobHandler(s services.JobService) *JobHandler {
	return &JobHandler{svc: s}
}

func (h *JobHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "HTTP Method invalid", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("GET request received at /jobs")

	jobs, err := h.svc.FindAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Found %d jobs", len(jobs))

	_ = json.NewEncoder(w).Encode(jobs)
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
		log.Printf("Job already on DB, updated successfully")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "Job updated successfully.",
		})
	case services.OutcomeUnchanged:
		log.Printf("Job already on DB, unchanged")
		w.WriteHeader(http.StatusNotModified)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "Job already on DB, unchanged.",
		})
	default:
		http.Error(w, "Unknown error", http.StatusInternalServerError)
	}
}

func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "HTTP Method invalid", http.StatusMethodNotAllowed)
		return
	}

	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		log.Println(err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	log.Printf("PUT request received at /jobs with payload: %v", job)

	err := h.svc.UpdateOnlyIfURLExists(r.Context(), job)

	if err != nil {
		if errors.Is(err, services.ErrJobNotFound) {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "Job updated successfully",
	})
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "HTTP Method invalid", http.StatusMethodNotAllowed)
		return
	}

	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		log.Println(err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	log.Printf("DELETE request received at /jobs with payload: %v", job)

	err := h.svc.DeleteOnlyIfURLExists(r.Context(), job)

	if err != nil {
		if errors.Is(err, services.ErrJobNotFound) {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "Job deleted successfully",
	})
}
