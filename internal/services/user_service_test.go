package services

import (
	"context"
	"errors"
	"jboard-go-crud/internal/models"
	"jboard-go-crud/internal/models/enums"
	"testing"
)

type mockUserRepository struct {
	createFunc         func(ctx context.Context, user models.User) error
	findByIDFunc       func(ctx context.Context, id string) (models.User, bool, error)
	findByUsernameFunc func(ctx context.Context, username string) (models.User, bool, error)
	updateByIDFunc     func(ctx context.Context, id string, user models.User) error
	deleteByIDFunc     func(ctx context.Context, id string) error
}

func (m *mockUserRepository) Create(ctx context.Context, user models.User) error {
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (models.User, bool, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockUserRepository) FindByUsername(ctx context.Context, username string) (models.User, bool, error) {
	return m.findByUsernameFunc(ctx, username)
}

func (m *mockUserRepository) UpdateByID(ctx context.Context, id string, user models.User) error {
	return m.updateByIDFunc(ctx, id, user)
}

func (m *mockUserRepository) DeleteByID(ctx context.Context, id string) error {
	return m.deleteByIDFunc(ctx, id)
}

func TestNewUserService(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}
}

func TestUserService_CreateUser_Success(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, nil
		},
		createFunc: func(ctx context.Context, user models.User) error {
			return nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "testuser", "password123", enums.Free)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUserService_CreateUser_EmptyUsername(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "", "password123", enums.Free)

	if err == nil {
		t.Error("Expected error for empty username, got nil")
	}

	if err.Error() != "username cannot be empty" {
		t.Errorf("Expected 'username cannot be empty', got %v", err)
	}
}

func TestUserService_CreateUser_WhitespaceUsername(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "   ", "password123", enums.Free)

	if err == nil {
		t.Error("Expected error for whitespace username, got nil")
	}

	if err.Error() != "username cannot be empty" {
		t.Errorf("Expected 'username cannot be empty', got %v", err)
	}
}

func TestUserService_CreateUser_EmptyPassword(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "testuser", "", enums.Free)

	if err == nil {
		t.Error("Expected error for empty password, got nil")
	}

	if err.Error() != "password cannot be empty" {
		t.Errorf("Expected 'password cannot be empty', got %v", err)
	}
}

func TestUserService_CreateUser_WhitespacePassword(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "testuser", "   ", enums.Free)

	if err == nil {
		t.Error("Expected error for whitespace password, got nil")
	}

	if err.Error() != "password cannot be empty" {
		t.Errorf("Expected 'password cannot be empty', got %v", err)
	}
}

func TestUserService_CreateUser_InvalidRole(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "testuser", "password123", "INVALID")

	if err == nil {
		t.Error("Expected error for invalid role, got nil")
	}

	if err.Error() != "invalid role" {
		t.Errorf("Expected 'invalid role', got %v", err)
	}
}

func TestUserService_CreateUser_UsernameAlreadyExists(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{ID: "existing-id", Username: username}, true, nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "existinguser", "password123", enums.Free)

	if err == nil {
		t.Error("Expected error for existing username, got nil")
	}

	if err.Error() != "username already exists" {
		t.Errorf("Expected 'username already exists', got %v", err)
	}
}

func TestUserService_CreateUser_RepositoryFindError(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, errors.New("database error")
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "testuser", "password123", enums.Free)

	if err == nil {
		t.Error("Expected error from repository, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected 'database error', got %v", err)
	}
}

func TestUserService_CreateUser_RepositoryCreateError(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, nil
		},
		createFunc: func(ctx context.Context, user models.User) error {
			return errors.New("create failed")
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "testuser", "password123", enums.Free)

	if err == nil {
		t.Error("Expected error from repository create, got nil")
	}

	if err.Error() != "create failed" {
		t.Errorf("Expected 'create failed', got %v", err)
	}
}

func TestUserService_CreateUser_PremiumRole(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, nil
		},
		createFunc: func(ctx context.Context, user models.User) error {
			return nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.CreateUser(ctx, "premiumuser", "password123", enums.Premium)

	if err != nil {
		t.Errorf("Expected no error for premium role, got %v", err)
	}
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	expectedUser := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	mockRepo := &mockUserRepository{
		findByIDFunc: func(ctx context.Context, id string) (models.User, bool, error) {
			return expectedUser, true, nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	user, err := service.GetUserByID(ctx, "test-id")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID != expectedUser.ID {
		t.Errorf("Expected user ID %s, got %s", expectedUser.ID, user.ID)
	}
}

func TestUserService_GetUserByID_EmptyID(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.GetUserByID(ctx, "")

	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}

	if err.Error() != "user ID cannot be empty" {
		t.Errorf("Expected 'user ID cannot be empty', got %v", err)
	}
}

func TestUserService_GetUserByID_WhitespaceID(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.GetUserByID(ctx, "   ")

	if err == nil {
		t.Error("Expected error for whitespace ID, got nil")
	}

	if err.Error() != "user ID cannot be empty" {
		t.Errorf("Expected 'user ID cannot be empty', got %v", err)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByIDFunc: func(ctx context.Context, id string) (models.User, bool, error) {
			return models.User{}, false, nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.GetUserByID(ctx, "nonexistent")

	if err == nil {
		t.Error("Expected error for user not found, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found', got %v", err)
	}
}

func TestUserService_GetUserByID_RepositoryError(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByIDFunc: func(ctx context.Context, id string) (models.User, bool, error) {
			return models.User{}, false, errors.New("repository error")
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.GetUserByID(ctx, "test-id")

	if err == nil {
		t.Error("Expected repository error, got nil")
	}

	if err.Error() != "repository error" {
		t.Errorf("Expected 'repository error', got %v", err)
	}
}

func TestUserService_GetUserByUsername_Success(t *testing.T) {
	expectedUser := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return expectedUser, true, nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	user, err := service.GetUserByUsername(ctx, "testuser")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.Username != expectedUser.Username {
		t.Errorf("Expected username %s, got %s", expectedUser.Username, user.Username)
	}
}

func TestUserService_GetUserByUsername_EmptyUsername(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.GetUserByUsername(ctx, "")

	if err == nil {
		t.Error("Expected error for empty username, got nil")
	}

	if err.Error() != "username cannot be empty" {
		t.Errorf("Expected 'username cannot be empty', got %v", err)
	}
}

func TestUserService_GetUserByUsername_NotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.GetUserByUsername(ctx, "nonexistent")

	if err == nil {
		t.Error("Expected error for user not found, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found', got %v", err)
	}
}

func TestUserService_GetUserByUsername_RepositoryError(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, errors.New("repository error")
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.GetUserByUsername(ctx, "testuser")

	if err == nil {
		t.Error("Expected repository error, got nil")
	}

	if err.Error() != "repository error" {
		t.Errorf("Expected 'repository error', got %v", err)
	}
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	existingUser := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "oldpass",
		Role:     enums.Free,
	}

	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return existingUser, true, nil
		},
		updateByIDFunc: func(ctx context.Context, id string, user models.User) error {
			return nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	updatedUser, err := service.UpdateUser(ctx, "testuser", "newpass", enums.Premium)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if updatedUser.Role != enums.Premium {
		t.Errorf("Expected role %s, got %s", enums.Premium, updatedUser.Role)
	}

	if updatedUser.Password != "" {
		t.Errorf("Expected password to be cleared in response, got %s", updatedUser.Password)
	}
}

func TestUserService_UpdateUser_EmptyUsername(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.UpdateUser(ctx, "", "newpass", enums.Premium)

	if err == nil {
		t.Error("Expected error for empty username, got nil")
	}

	if err.Error() != "username cannot be empty" {
		t.Errorf("Expected 'username cannot be empty', got %v", err)
	}
}

func TestUserService_UpdateUser_InvalidRole(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.UpdateUser(ctx, "testuser", "newpass", "INVALID")

	if err == nil {
		t.Error("Expected error for invalid role, got nil")
	}

	if err.Error() != "invalid role" {
		t.Errorf("Expected 'invalid role', got %v", err)
	}
}

func TestUserService_UpdateUser_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.UpdateUser(ctx, "nonexistent", "newpass", enums.Premium)

	if err == nil {
		t.Error("Expected error for user not found, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found', got %v", err)
	}
}

func TestUserService_UpdateUser_EmptyPassword(t *testing.T) {
	existingUser := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "oldpass",
		Role:     enums.Free,
	}

	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return existingUser, true, nil
		},
		updateByIDFunc: func(ctx context.Context, id string, user models.User) error {
			if user.Password != existingUser.Password {
				t.Errorf("Expected password to remain unchanged, got %s", user.Password)
			}
			return nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.UpdateUser(ctx, "testuser", "", enums.Premium)

	if err != nil {
		t.Errorf("Expected no error for empty password, got %v", err)
	}
}

func TestUserService_UpdateUser_RepositoryFindError(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, errors.New("find error")
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.UpdateUser(ctx, "testuser", "newpass", enums.Premium)

	if err == nil {
		t.Error("Expected repository find error, got nil")
	}

	if err.Error() != "find error" {
		t.Errorf("Expected 'find error', got %v", err)
	}
}

func TestUserService_UpdateUser_RepositoryUpdateError(t *testing.T) {
	existingUser := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "oldpass",
		Role:     enums.Free,
	}

	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return existingUser, true, nil
		},
		updateByIDFunc: func(ctx context.Context, id string, user models.User) error {
			return errors.New("update error")
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	_, err := service.UpdateUser(ctx, "testuser", "newpass", enums.Premium)

	if err == nil {
		t.Error("Expected repository update error, got nil")
	}

	if err.Error() != "update error" {
		t.Errorf("Expected 'update error', got %v", err)
	}
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	existingUser := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return existingUser, true, nil
		},
		deleteByIDFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.DeleteUser(ctx, "testuser")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUserService_DeleteUser_EmptyUsername(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.DeleteUser(ctx, "")

	if err == nil {
		t.Error("Expected error for empty username, got nil")
	}

	if err.Error() != "username cannot be empty" {
		t.Errorf("Expected 'username cannot be empty', got %v", err)
	}
}

func TestUserService_DeleteUser_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, nil
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.DeleteUser(ctx, "nonexistent")

	if err == nil {
		t.Error("Expected error for user not found, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found', got %v", err)
	}
}

func TestUserService_DeleteUser_RepositoryFindError(t *testing.T) {
	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return models.User{}, false, errors.New("find error")
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.DeleteUser(ctx, "testuser")

	if err == nil {
		t.Error("Expected repository find error, got nil")
	}

	if err.Error() != "find error" {
		t.Errorf("Expected 'find error', got %v", err)
	}
}

func TestUserService_DeleteUser_RepositoryDeleteError(t *testing.T) {
	existingUser := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	mockRepo := &mockUserRepository{
		findByUsernameFunc: func(ctx context.Context, username string) (models.User, bool, error) {
			return existingUser, true, nil
		},
		deleteByIDFunc: func(ctx context.Context, id string) error {
			return errors.New("delete error")
		},
	}

	service := NewUserService(mockRepo)
	ctx := context.Background()

	err := service.DeleteUser(ctx, "testuser")

	if err == nil {
		t.Error("Expected repository delete error, got nil")
	}

	if err.Error() != "delete error" {
		t.Errorf("Expected 'delete error', got %v", err)
	}
}
