package controllers

import (
	"encoding/json"
	"jboard-go-crud/internal/models/enums"
	"jboard-go-crud/internal/services"
	"log"
	"net/http"
	"strings"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	log.Printf("Creating new UserHandler")
	return &UserHandler{
		userService: userService,
	}
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler CreateUser called")

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode create user request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role := enums.RoleEnum(req.Role)
	if !role.IsValid() {
		log.Printf("Invalid role provided: %s", req.Role)
		http.Error(w, "Invalid role. Must be Free or Premium", http.StatusBadRequest)
		return
	}

	err := h.userService.CreateUser(r.Context(), req.Username, req.Password, role)
	if err != nil {
		log.Printf("Service error in CreateUser: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "cannot be empty") || strings.Contains(err.Error(), "invalid") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler GetUser called")

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("ID query parameter is missing")
		http.Error(w, "ID query parameter is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		log.Printf("Service error in GetUser: %v", err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "cannot be empty") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler GetUserByUsername called")

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Printf("Username query parameter is missing")
		http.Error(w, "Username query parameter is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByUsername(r.Context(), username)
	if err != nil {
		log.Printf("Service error in GetUserByUsername: %v", err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "cannot be empty") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if username := r.URL.Query().Get("username"); username != "" {
		h.GetUserByUsername(w, r)
		return
	}

	if id := r.URL.Query().Get("id"); id != "" {
		h.GetUser(w, r)
		return
	}

	log.Printf("Neither id nor username query parameter provided")
	http.Error(w, "Either 'id' or 'username' query parameter is required", http.StatusBadRequest)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler UpdateUser called")

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode update user request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	role := enums.RoleEnum(req.Role)
	if !role.IsValid() {
		log.Printf("Invalid role provided: %s", req.Role)
		http.Error(w, "Invalid role. Must be Free or Premium", http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), req.Username, req.Password, role)
	if err != nil {
		log.Printf("Service error in UpdateUser: %v", err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "cannot be empty") || strings.Contains(err.Error(), "invalid") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler DeleteUser called")

	username := r.URL.Query().Get("username")

	err := h.userService.DeleteUser(r.Context(), username)
	if err != nil {
		log.Printf("Service error in DeleteUser: %v", err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "cannot be empty") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
