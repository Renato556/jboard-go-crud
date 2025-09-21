package routers

import (
	"context"
	"jboard-go-crud/src/controllers"
	"jboard-go-crud/src/models"
	"jboard-go-crud/src/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockJobService struct{}

func (m *mockJobService) CreateOrUpdate(_ context.Context, _ models.Job) (services.UpsertOutcome, error) {
	return services.OutcomeCreated, nil
}

func (m *mockJobService) FindAll(_ context.Context) ([]models.Job, error) {
	return []models.Job{
		{ID: "test-1", Title: "Job 1"},
		{ID: "test-2", Title: "Job 2"},
	}, nil
}

func TestNewJobsController(t *testing.T) {
	mockService := &mockJobService{}
	jobHandler := controllers.NewJobHandler(mockService)

	handler := NewJobsController(jobHandler)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}

func TestNewJobsController_PostRoute(t *testing.T) {
	mockService := &mockJobService{}
	jobHandler := controllers.NewJobHandler(mockService)

	handler := NewJobsController(jobHandler)

	req := httptest.NewRequest(http.MethodPost, "/v1/jobs", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Error("POST /v1/jobs route should be registered")
	}
}

func TestNewJobsController_GetRouteNotAllowed(t *testing.T) {
	mockService := &mockJobService{}
	jobHandler := controllers.NewJobHandler(mockService)

	handler := NewJobsController(jobHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/jobs/invalid", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for invalid GET request, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNewJobsController_PutRouteNotFound(t *testing.T) {
	mockService := &mockJobService{}
	jobHandler := controllers.NewJobHandler(mockService)

	handler := NewJobsController(jobHandler)

	req := httptest.NewRequest(http.MethodPut, "/v1/jobs", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for PUT request, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestNewJobsController_DeleteRouteNotFound(t *testing.T) {
	mockService := &mockJobService{}
	jobHandler := controllers.NewJobHandler(mockService)

	handler := NewJobsController(jobHandler)

	req := httptest.NewRequest(http.MethodDelete, "/v1/jobs", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for DELETE request, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestNewJobsController_InvalidRoute(t *testing.T) {
	mockService := &mockJobService{}
	jobHandler := controllers.NewJobHandler(mockService)

	handler := NewJobsController(jobHandler)

	req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for invalid route, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHealthCheck(t *testing.T) {
	mockService := &mockJobService{}
	jobHandler := controllers.NewJobHandler(mockService)

	handler := NewJobsController(jobHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d for health check, got %d", http.StatusOK, rr.Code)
	}

	expectedBody := `{"status":"healthy","service":"jboard-go-crud"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, rr.Body.String())
	}

	expectedContentType := "application/json"
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, rr.Header().Get("Content-Type"))
	}
}

func TestGetAllJobsRoute(t *testing.T) {
	mockService := &mockJobService{}
	jobHandler := controllers.NewJobHandler(mockService)

	handler := NewJobsController(jobHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/jobs", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d for GET /v1/jobs, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}
