package controllers

import (
	"encoding/json"
	"jboard-go-crud/internal/models"
	"jboard-go-crud/internal/services"
	"log"
	"net/http"
	"strings"
)

type SkillHandler struct {
	skillService services.SkillService
}

func NewSkillHandler(skillService services.SkillService) *SkillHandler {
	log.Printf("Creating new SkillHandler")
	return &SkillHandler{
		skillService: skillService,
	}
}

func (h *SkillHandler) GetAllSkills(w http.ResponseWriter, r *http.Request) {
	log.Printf("Controller GetAllSkills called")

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Printf("ERROR: Missing username query parameter in GetAllSkills")
		http.Error(w, "username query parameter is required", http.StatusBadRequest)
		return
	}
	log.Printf("Controller GetAllSkills processing request for username: %s", username)

	skill, err := h.skillService.GetAllSkills(r.Context(), username)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("User not found in GetAllSkills for username: %s", username)
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			log.Printf("ERROR: Service error in GetAllSkills for username %s: %v", username, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully retrieved skills for username: %s, returning response", username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skill)
}

func (h *SkillHandler) AddSkill(w http.ResponseWriter, r *http.Request) {
	log.Printf("Controller AddSkill called")

	var skillRequest models.SkillRequest
	if err := json.NewDecoder(r.Body).Decode(&skillRequest); err != nil {
		log.Printf("ERROR: Invalid JSON in AddSkill request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("Controller AddSkill processing request for username: %s, skill: %s", skillRequest.Username, skillRequest.Skill)

	if err := h.skillService.AddSkill(r.Context(), skillRequest); err != nil {
		if strings.Contains(err.Error(), "required") {
			log.Printf("Validation error in AddSkill: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			log.Printf("ERROR: Service error in AddSkill for username %s: %v", skillRequest.Username, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully added skill for username: %s, returning response", skillRequest.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Skill added successfully"})
}

func (h *SkillHandler) RemoveSkill(w http.ResponseWriter, r *http.Request) {
	log.Printf("Controller RemoveSkill called")

	var skillRequest models.SkillRequest
	if err := json.NewDecoder(r.Body).Decode(&skillRequest); err != nil {
		log.Printf("ERROR: Invalid JSON in RemoveSkill request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("Controller RemoveSkill processing request for username: %s, skill: %s", skillRequest.Username, skillRequest.Skill)

	if err := h.skillService.RemoveSkill(r.Context(), skillRequest); err != nil {
		if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "not found") {
			log.Printf("Validation or not found error in RemoveSkill: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			log.Printf("ERROR: Service error in RemoveSkill for username %s: %v", skillRequest.Username, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully removed skill for username: %s, returning response", skillRequest.Username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Skill removed successfully"})
}

func (h *SkillHandler) DeleteUserSkills(w http.ResponseWriter, r *http.Request) {
	log.Printf("Controller DeleteUserSkills called")

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Printf("ERROR: Missing username query parameter in DeleteUserSkills")
		http.Error(w, "username query parameter is required", http.StatusBadRequest)
		return
	}
	log.Printf("Controller DeleteUserSkills processing request for username: %s", username)

	if err := h.skillService.DeleteUserSkills(r.Context(), username); err != nil {
		if strings.Contains(err.Error(), "user not found") {
			log.Printf("User not found error in DeleteUserSkills: %v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if strings.Contains(err.Error(), "required") {
			log.Printf("Validation error in DeleteUserSkills: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			log.Printf("ERROR: Service error in DeleteUserSkills for username %s: %v", username, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully deleted skills for username: %s, returning response", username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User skills deleted successfully"})
}
