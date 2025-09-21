package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"jboard-go-crud/src/models"
	"jboard-go-crud/src/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockJobService struct {
	createOrUpdateFunc func(ctx context.Context, job models.Job) (services.UpsertOutcome, error)
	findAllFunc        func(ctx context.Context) ([]models.Job, error)
}

func (m *mockJobService) CreateOrUpdate(ctx context.Context, job models.Job) (services.UpsertOutcome, error) {
	return m.createOrUpdateFunc(ctx, job)
}

func (m *mockJobService) FindAll(ctx context.Context) ([]models.Job, error) {
	return m.findAllFunc(ctx)
}

func TestNewJobHandler(t *testing.T) {
	mockService := &mockJobService{}
	handler := NewJobHandler(mockService)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}

func TestJobHandler_CreateJob_Success_Created(t *testing.T) {
	job := models.Job{
		ID:             "test-id",
		Title:          "Test Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
		IsBrazilianFriendly: models.BrazilianFriendly{
			IsFriendly: true,
			Reason:     "Remote",
		},
	}

	mockService := &mockJobService{
		createOrUpdateFunc: func(ctx context.Context, job models.Job) (services.UpsertOutcome, error) {
			return services.OutcomeCreated, nil
		},
	}

	handler := NewJobHandler(mockService)

	jobJSON, _ := json.Marshal(job)
	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(jobJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateJob(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	if response["message"] != "Job created successfully." {
		t.Errorf("Expected message 'Job created successfully.', got %v", response["message"])
	}
}

func TestJobHandler_CreateJob_Success_Updated(t *testing.T) {
	job := models.Job{
		ID:             "test-id",
		Title:          "Test Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
		IsBrazilianFriendly: models.BrazilianFriendly{
			IsFriendly: true,
			Reason:     "Remote",
		},
	}

	mockService := &mockJobService{
		createOrUpdateFunc: func(ctx context.Context, job models.Job) (services.UpsertOutcome, error) {
			return services.OutcomeUpdated, nil
		},
	}

	handler := NewJobHandler(mockService)

	jobJSON, _ := json.Marshal(job)
	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(jobJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateJob(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	if response["message"] != "Job already exists, updated with new information and extended expiration." {
		t.Errorf("Expected updated message, got %v", response["message"])
	}
}

func TestJobHandler_CreateJob_InvalidMethod(t *testing.T) {
	mockService := &mockJobService{}
	handler := NewJobHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	rr := httptest.NewRecorder()

	handler.CreateJob(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestJobHandler_CreateJob_InvalidPayload(t *testing.T) {
	mockService := &mockJobService{}
	handler := NewJobHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateJob(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestJobHandler_CreateJob_ServiceError(t *testing.T) {
	job := models.Job{
		ID:             "test-id",
		Title:          "Test Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
		IsBrazilianFriendly: models.BrazilianFriendly{
			IsFriendly: true,
			Reason:     "Remote",
		},
	}

	mockService := &mockJobService{
		createOrUpdateFunc: func(ctx context.Context, job models.Job) (services.UpsertOutcome, error) {
			return 0, errors.New("service error")
		},
	}

	handler := NewJobHandler(mockService)

	jobJSON, _ := json.Marshal(job)
	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(jobJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateJob(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestJobHandler_CreateJob_UnknownOutcome(t *testing.T) {
	job := models.Job{
		ID:             "test-id",
		Title:          "Test Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
		IsBrazilianFriendly: models.BrazilianFriendly{
			IsFriendly: true,
			Reason:     "Remote",
		},
	}

	mockService := &mockJobService{
		createOrUpdateFunc: func(ctx context.Context, job models.Job) (services.UpsertOutcome, error) {
			return 999, nil
		},
	}

	handler := NewJobHandler(mockService)

	jobJSON, _ := json.Marshal(job)
	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(jobJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateJob(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestJobHandler_GetAllJobs_Success(t *testing.T) {
	expectedJobs := []models.Job{
		{
			ID:             "test-id-1",
			Title:          "Job 1",
			Company:        "Company 1",
			Url:            "https://test1.com",
			SeniorityLevel: "Senior",
			Field:          "Engineering",
		},
		{
			ID:             "test-id-2",
			Title:          "Job 2",
			Company:        "Company 2",
			Url:            "https://test2.com",
			SeniorityLevel: "Junior",
			Field:          "Design",
		},
	}

	mockService := &mockJobService{
		findAllFunc: func(ctx context.Context) ([]models.Job, error) {
			return expectedJobs, nil
		},
	}

	handler := NewJobHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/v1/jobs", nil)
	rr := httptest.NewRecorder()

	handler.GetAllJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var response []models.Job
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(response))
	}

	if response[0].ID != "test-id-1" {
		t.Errorf("Expected first job ID to be 'test-id-1', got %s", response[0].ID)
	}
}

func TestJobHandler_GetAllJobs_InvalidMethod(t *testing.T) {
	mockService := &mockJobService{}
	handler := NewJobHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/v1/jobs", nil)
	rr := httptest.NewRecorder()

	handler.GetAllJobs(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestJobHandler_GetAllJobs_ServiceError(t *testing.T) {
	mockService := &mockJobService{
		findAllFunc: func(ctx context.Context) ([]models.Job, error) {
			return nil, errors.New("service error")
		},
	}

	handler := NewJobHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/v1/jobs", nil)
	rr := httptest.NewRecorder()

	handler.GetAllJobs(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
