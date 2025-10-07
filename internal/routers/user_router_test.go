package routers

import (
	"context"
	"jboard-go-crud/internal/controllers"
	"jboard-go-crud/internal/models"
	"jboard-go-crud/internal/models/enums"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockUserService struct{}

func (m *mockUserService) CreateUser(_ context.Context, _, _ string, _ enums.RoleEnum) error {
	return nil
}

func (m *mockUserService) GetUserByID(_ context.Context, _ string) (models.User, error) {
	testID, _ := primitive.ObjectIDFromHex("68e462f868efefe99e226a8b")
	return models.User{ID: testID, Username: "testuser", Role: enums.Free}, nil
}

func (m *mockUserService) GetUserByUsername(_ context.Context, _ string) (models.User, error) {
	testID, _ := primitive.ObjectIDFromHex("68e462f868efefe99e226a8b")
	return models.User{ID: testID, Username: "testuser", Role: enums.Free}, nil
}

func (m *mockUserService) UpdateUser(_ context.Context, _, _ string, _ enums.RoleEnum) (models.User, error) {
	testID, _ := primitive.ObjectIDFromHex("68e462f868efefe99e226a8b")
	return models.User{ID: testID, Username: "testuser", Role: enums.Premium}, nil
}

func (m *mockUserService) DeleteUser(_ context.Context, _ string) error {
	return nil
}

func TestNewUsersController(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}

func TestNewUsersController_PostRoute(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodPost, "/v1/users", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Error("POST /v1/users route should be registered")
	}
}

func TestNewUsersController_GetRoute(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Error("GET /v1/users route should be registered")
	}
}

func TestNewUsersController_PutRoute(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodPut, "/v1/users", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Error("PUT /v1/users route should be registered")
	}
}

func TestNewUsersController_DeleteRoute(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodDelete, "/v1/users", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Error("DELETE /v1/users route should be registered")
	}
}

func TestNewUsersController_InvalidRoute(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/users/invalid", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for invalid route, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNewUsersController_PatchMethodNotAllowed(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodPatch, "/v1/users", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for PATCH request, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestNewUsersController_HeadMethodNotAllowed(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodHead, "/v1/users?id=test-id", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d for HEAD request with valid parameters, got %d", http.StatusOK, rr.Code)
	}
}

func TestNewUsersController_OptionsMethodNotAllowed(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodOptions, "/v1/users", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for OPTIONS request, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestNewUsersController_WithNilHandler(t *testing.T) {
	handler := NewUsersController(nil)

	if handler == nil {
		t.Error("Expected handler to be created even with nil userHandler, got nil")
	}
}

func TestNewUsersController_RoutePathSensitivity(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodGet, "/V1/USERS", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for case-sensitive route mismatch, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNewUsersController_RouteWithQueryParameters(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/users?id=123", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusNotFound {
		t.Error("GET /v1/users with query parameters should be routed correctly")
	}
}

func TestNewUsersController_RouteWithTrailingSlash(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodGet, "/v1/users/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for route with trailing slash, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNewUsersController_EmptyPath(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for empty path, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNewUsersController_WrongVersionPath(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	req := httptest.NewRequest(http.MethodGet, "/v2/users", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for wrong version path, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestNewUsersController_MultipleRequestsToSameRoute(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := controllers.NewUserHandler(mockService)

	handler := NewUsersController(userHandler)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code == http.StatusNotFound {
			t.Errorf("GET /v1/users route should handle multiple requests, failed on request %d", i+1)
		}
	}
}
