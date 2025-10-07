package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"jboard-go-crud/internal/models"
	"jboard-go-crud/internal/models/enums"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserService struct {
	createUserFunc        func(ctx context.Context, username, password string, role enums.RoleEnum) error
	getUserByIDFunc       func(ctx context.Context, id string) (models.User, error)
	getUserByUsernameFunc func(ctx context.Context, username string) (models.User, error)
	updateUserFunc        func(ctx context.Context, username string, password string, role enums.RoleEnum) (models.User, error)
	deleteUserFunc        func(ctx context.Context, username string) error
}

func (m *mockUserService) CreateUser(ctx context.Context, username, password string, role enums.RoleEnum) error {
	return m.createUserFunc(ctx, username, password, role)
}

func (m *mockUserService) GetUserByID(ctx context.Context, id string) (models.User, error) {
	return m.getUserByIDFunc(ctx, id)
}

func (m *mockUserService) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return m.getUserByUsernameFunc(ctx, username)
}

func (m *mockUserService) UpdateUser(ctx context.Context, username string, password string, role enums.RoleEnum) (models.User, error) {
	return m.updateUserFunc(ctx, username, password, role)
}

func (m *mockUserService) DeleteUser(ctx context.Context, username string) error {
	return m.deleteUserFunc(ctx, username)
}

func TestNewUserHandler(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	mockService := &mockUserService{
		createUserFunc: func(ctx context.Context, username, password string, role enums.RoleEnum) error {
			return nil
		},
	}

	handler := NewUserHandler(mockService)

	reqBody := CreateUserRequest{
		Username: "testuser",
		Password: "testpass",
		Role:     "FREE",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_CreateUser_InvalidRole(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	reqBody := CreateUserRequest{
		Username: "testuser",
		Password: "testpass",
		Role:     "INVALID",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_CreateUser_UserAlreadyExists(t *testing.T) {
	mockService := &mockUserService{
		createUserFunc: func(ctx context.Context, username, password string, role enums.RoleEnum) error {
			return errors.New("user already exists")
		},
	}

	handler := NewUserHandler(mockService)

	reqBody := CreateUserRequest{
		Username: "existinguser",
		Password: "testpass",
		Role:     "FREE",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, rr.Code)
	}
}

func TestUserHandler_CreateUser_EmptyUsernameValidationError(t *testing.T) {
	mockService := &mockUserService{
		createUserFunc: func(ctx context.Context, username, password string, role enums.RoleEnum) error {
			return errors.New("username cannot be empty")
		},
	}

	handler := NewUserHandler(mockService)

	reqBody := CreateUserRequest{
		Username: "",
		Password: "testpass",
		Role:     "FREE",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_CreateUser_InternalServerError(t *testing.T) {
	mockService := &mockUserService{
		createUserFunc: func(ctx context.Context, username, password string, role enums.RoleEnum) error {
			return errors.New("database connection failed")
		},
	}

	handler := NewUserHandler(mockService)

	reqBody := CreateUserRequest{
		Username: "testuser",
		Password: "testpass",
		Role:     "FREE",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	testID := primitive.NewObjectID()
	expectedUser := models.User{
		ID:       testID,
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	mockService := &mockUserService{
		getUserByIDFunc: func(ctx context.Context, id string) (models.User, error) {
			return expectedUser, nil
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users?id="+testID.Hex(), nil)
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var responseUser models.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	if responseUser.ID != expectedUser.ID {
		t.Errorf("Expected user ID %s, got %s", expectedUser.ID.Hex(), responseUser.ID.Hex())
	}
}

func TestUserHandler_GetUser_MissingIDParameter(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_GetUser_UserNotFound(t *testing.T) {
	mockService := &mockUserService{
		getUserByIDFunc: func(ctx context.Context, id string) (models.User, error) {
			return models.User{}, errors.New("user not found")
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users?id=nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestUserHandler_GetUser_EmptyIDValidationError(t *testing.T) {
	mockService := &mockUserService{
		getUserByIDFunc: func(ctx context.Context, id string) (models.User, error) {
			return models.User{}, errors.New("id cannot be empty")
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users?id=test-id", nil)
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_GetUser_InternalServerError(t *testing.T) {
	mockService := &mockUserService{
		getUserByIDFunc: func(ctx context.Context, id string) (models.User, error) {
			return models.User{}, errors.New("database connection failed")
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users?id=test-id", nil)
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestUserHandler_GetUserByUsername_Success(t *testing.T) {
	testID := primitive.NewObjectID()
	expectedUser := models.User{
		ID:       testID,
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	mockService := &mockUserService{
		getUserByUsernameFunc: func(ctx context.Context, username string) (models.User, error) {
			return expectedUser, nil
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users?username=testuser", nil)
	rr := httptest.NewRecorder()

	handler.GetUserByUsername(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var responseUser models.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	if responseUser.Username != expectedUser.Username {
		t.Errorf("Expected username %s, got %s", expectedUser.Username, responseUser.Username)
	}
}

func TestUserHandler_GetUserByUsername_MissingUsernameParameter(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler.GetUserByUsername(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_GetUserByUsername_UserNotFound(t *testing.T) {
	mockService := &mockUserService{
		getUserByUsernameFunc: func(ctx context.Context, username string) (models.User, error) {
			return models.User{}, errors.New("user not found")
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users?username=nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.GetUserByUsername(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestUserHandler_GetUserHandler_WithUsernameParameter(t *testing.T) {
	testID := primitive.NewObjectID()
	expectedUser := models.User{
		ID:       testID,
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	mockService := &mockUserService{
		getUserByUsernameFunc: func(ctx context.Context, username string) (models.User, error) {
			return expectedUser, nil
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users?username=testuser", nil)
	rr := httptest.NewRecorder()

	handler.GetUserHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestUserHandler_GetUserHandler_WithIDParameter(t *testing.T) {
	testID := primitive.NewObjectID()
	expectedUser := models.User{
		ID:       testID,
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	mockService := &mockUserService{
		getUserByIDFunc: func(ctx context.Context, id string) (models.User, error) {
			return expectedUser, nil
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users?id="+testID.Hex(), nil)
	rr := httptest.NewRecorder()

	handler.GetUserHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestUserHandler_GetUserHandler_MissingBothParameters(t *testing.T) {
	mockService := &mockUserService{}
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler.GetUserHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	testID := primitive.NewObjectID()
	updatedUser := models.User{
		ID:       testID,
		Username: "testuser",
		Password: "newhashedpass",
		Role:     enums.Premium,
	}

	mockService := &mockUserService{
		updateUserFunc: func(ctx context.Context, username string, password string, role enums.RoleEnum) (models.User, error) {
			return updatedUser, nil
		},
	}

	handler := NewUserHandler(mockService)

	reqBody := UpdateUserRequest{
		Username: "testuser",
		Password: "newpass",
		Role:     "PREMIUM",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.UpdateUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var responseUser models.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	if responseUser.Role != enums.Premium {
		t.Errorf("Expected role %s, got %s", enums.Premium, responseUser.Role)
	}
}

// func TestUserHandler_UpdateUser_MissingUsername(t *testing.T) {
// 	mockService := &mockUserService{}
// 	handler := NewUserHandler(mockService)

// 	reqBody := UpdateUserRequest{
// 		Password: "newpass",
// 		Role:     "PREMIUM",
// 	}
// 	reqJSON, _ := json.Marshal(reqBody)
// 	req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(reqJSON))
// 	req.Header.Set("Content-Type", "application/json")

// 	rr := httptest.NewRecorder()

// 	handler.UpdateUser(rr, req)

// 	if rr.Code != http.StatusBadRequest {
// 		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
// 	}
// }

// func TestUserHandler_UpdateUser_UserNotFound(t *testing.T) {
// 	mockService := &mockUserService{
// 		updateUserFunc: func(ctx context.Context, username string, password string, role enums.RoleEnum) (models.User, error) {
// 			return models.User{}, errors.New("user not found")
// 		},
// 	}

// 	handler := NewUserHandler(mockService)

// 	reqBody := UpdateUserRequest{
// 		Username: "nonexistent",
// 		Password: "newpass",
// 		Role:     "PREMIUM",
// 	}
// 	reqJSON, _ := json.Marshal(reqBody)
// 	req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(reqJSON))
// 	req.Header.Set("Content-Type", "application/json")

// 	rr := httptest.NewRecorder()

// 	handler.UpdateUser(rr, req)

// 	if rr.Code != http.StatusNotFound {
// 		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
// 	}
// }

// func TestUserHandler_UpdateUser_InternalServerError(t *testing.T) {
// 	mockService := &mockUserService{
// 		updateUserFunc: func(ctx context.Context, username string, password string, role enums.RoleEnum) (models.User, error) {
// 			return models.User{}, errors.New("database connection failed")
// 		},
// 	}

// 	handler := NewUserHandler(mockService)

// 	reqBody := UpdateUserRequest{
// 		Username: "testuser",
// 		Password: "newpass",
// 		Role:     "PREMIUM",
// 	}
// 	reqJSON, _ := json.Marshal(reqBody)
// 	req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(reqJSON))
// 	req.Header.Set("Content-Type", "application/json")

// 	rr := httptest.NewRecorder()

// 	handler.UpdateUser(rr, req)

// 	if rr.Code != http.StatusInternalServerError {
// 		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
// 	}
// }

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	mockService := &mockUserService{
		deleteUserFunc: func(ctx context.Context, username string) error {
			return nil
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/users?username=testuser", nil)
	rr := httptest.NewRecorder()

	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestUserHandler_DeleteUser_MissingUsernameParameter(t *testing.T) {
	mockService := &mockUserService{
		deleteUserFunc: func(ctx context.Context, username string) error {
			return nil // This won't be called since the handler should return before calling the service
		},
	}
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	rr := httptest.NewRecorder()

	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_DeleteUser_UserNotFound(t *testing.T) {
	mockService := &mockUserService{
		deleteUserFunc: func(ctx context.Context, username string) error {
			return errors.New("user not found")
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/users?username=nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestUserHandler_DeleteUser_InternalServerError(t *testing.T) {
	mockService := &mockUserService{
		deleteUserFunc: func(ctx context.Context, username string) error {
			return errors.New("database connection failed")
		},
	}

	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/users?username=testuser", nil)
	rr := httptest.NewRecorder()

	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
