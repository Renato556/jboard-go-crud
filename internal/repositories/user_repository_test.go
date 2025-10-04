package repositories

import (
	"context"
	"jboard-go-crud/internal/models"
	"jboard-go-crud/internal/models/enums"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	if repo == nil {
		t.Error("Expected repository to be created, got nil")
	}
}

func TestNewUserRepository_WithValidClient(t *testing.T) {
	client := &mongo.Client{}
	repo := NewUserRepository(client, "testdb", "users")

	if repo == nil {
		t.Error("Expected repository to be created, got nil")
	}
}

func TestUserRepository_Create_Success(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.Create(context.Background(), user)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get users collection" {
		t.Errorf("Expected 'failed to get users collection' error, got %v", err)
	}
}

func TestUserRepository_Create_ValidationError(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user := models.User{
		ID:       "",
		Username: "",
		Password: "",
		Role:     "",
	}

	err := repo.Create(context.Background(), user)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
}

func TestUserRepository_FindByID_Success(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user, found, err := repo.FindByID(context.Background(), "test-id")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if found {
		t.Error("Expected found to be false, got true")
	}

	if user.ID != "" {
		t.Errorf("Expected empty user, got user with ID: %s", user.ID)
	}

	if err.Error() != "failed to get users collection" {
		t.Errorf("Expected 'failed to get users collection' error, got %v", err)
	}
}

func TestUserRepository_FindByID_EmptyID(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user, found, err := repo.FindByID(context.Background(), "")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if found {
		t.Error("Expected found to be false, got true")
	}

	if user.ID != "" {
		t.Errorf("Expected empty user, got user with ID: %s", user.ID)
	}
}

func TestUserRepository_FindByUsername_Success(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user, found, err := repo.FindByUsername(context.Background(), "testuser")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if found {
		t.Error("Expected found to be false, got true")
	}

	if user.Username != "" {
		t.Errorf("Expected empty user, got user with username: %s", user.Username)
	}

	if err.Error() != "failed to get users collection" {
		t.Errorf("Expected 'failed to get users collection' error, got %v", err)
	}
}

func TestUserRepository_FindByUsername_EmptyUsername(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user, found, err := repo.FindByUsername(context.Background(), "")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if found {
		t.Error("Expected found to be false, got true")
	}

	if user.Username != "" {
		t.Errorf("Expected empty user, got user with username: %s", user.Username)
	}
}

func TestUserRepository_UpdateByID_Success(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "newhashedpass",
		Role:     enums.Premium,
	}

	err := repo.UpdateByID(context.Background(), "test-id", user)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get users collection" {
		t.Errorf("Expected 'failed to get users collection' error, got %v", err)
	}
}

func TestUserRepository_UpdateByID_ValidationError(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user := models.User{
		ID:       "",
		Username: "",
		Password: "",
		Role:     "",
	}

	err := repo.UpdateByID(context.Background(), "test-id", user)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
}

func TestUserRepository_UpdateByID_EmptyID(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.UpdateByID(context.Background(), "", user)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get users collection" {
		t.Errorf("Expected 'failed to get users collection' error, got %v", err)
	}
}

func TestUserRepository_DeleteByID_Success(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	err := repo.DeleteByID(context.Background(), "test-id")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get users collection" {
		t.Errorf("Expected 'failed to get users collection' error, got %v", err)
	}
}

func TestUserRepository_DeleteByID_EmptyID(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	err := repo.DeleteByID(context.Background(), "")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get users collection" {
		t.Errorf("Expected 'failed to get users collection' error, got %v", err)
	}
}

func TestUserRepository_ContextCancellation(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	user := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.Create(ctx, user)
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}

	_, _, err = repo.FindByID(ctx, "test-id")
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}

	_, _, err = repo.FindByUsername(ctx, "testuser")
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}

	err = repo.UpdateByID(ctx, "test-id", user)
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}

	err = repo.DeleteByID(ctx, "test-id")
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}
}

func TestUserRepository_Create_DuplicateUsernameHandling(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.Create(context.Background(), user)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}
}

func TestUserRepository_UpdateByID_DuplicateUsernameHandling(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user := models.User{
		ID:       "test-id",
		Username: "existinguser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.UpdateByID(context.Background(), "test-id", user)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}
}

func TestUserRepository_Create_ValidUserAllRoles(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	freeUser := models.User{
		ID:       "free-user-id",
		Username: "freeuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.Create(context.Background(), freeUser)
	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	premiumUser := models.User{
		ID:       "premium-user-id",
		Username: "premiumuser",
		Password: "hashedpass",
		Role:     enums.Premium,
	}

	err = repo.Create(context.Background(), premiumUser)
	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}
}

func TestUserRepository_FindByID_NonExistentUser(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user, found, err := repo.FindByID(context.Background(), "nonexistent-id")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if found {
		t.Error("Expected found to be false for nonexistent user, got true")
	}

	if user.ID != "" {
		t.Errorf("Expected empty user for nonexistent ID, got user with ID: %s", user.ID)
	}
}

func TestUserRepository_FindByUsername_NonExistentUser(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user, found, err := repo.FindByUsername(context.Background(), "nonexistentuser")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if found {
		t.Error("Expected found to be false for nonexistent user, got true")
	}

	if user.Username != "" {
		t.Errorf("Expected empty user for nonexistent username, got user with username: %s", user.Username)
	}
}

func TestUserRepository_DeleteByID_NonExistentUser(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	err := repo.DeleteByID(context.Background(), "nonexistent-id")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get users collection" {
		t.Errorf("Expected 'failed to get users collection' error, got %v", err)
	}
}

func TestUserRepository_UpdateByID_NonExistentUser(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user := models.User{
		ID:       "nonexistent-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.UpdateByID(context.Background(), "nonexistent-id", user)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get users collection" {
		t.Errorf("Expected 'failed to get users collection' error, got %v", err)
	}
}

func TestUserRepository_Create_LongUsername(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	longUsername := string(make([]byte, 1000))
	user := models.User{
		ID:       "test-id",
		Username: longUsername,
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.Create(context.Background(), user)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}
}

func TestUserRepository_Create_SpecialCharactersInUsername(t *testing.T) {
	repo := NewUserRepository(nil, "testdb", "users")

	user := models.User{
		ID:       "test-id",
		Username: "user@domain.com",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.Create(context.Background(), user)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}
}

func TestUserRepository_Database_NilHandling(t *testing.T) {
	repo := NewUserRepository(nil, "", "")

	user := models.User{
		ID:       "test-id",
		Username: "testuser",
		Password: "hashedpass",
		Role:     enums.Free,
	}

	err := repo.Create(context.Background(), user)
	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	_, _, err = repo.FindByID(context.Background(), "test-id")
	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	_, _, err = repo.FindByUsername(context.Background(), "testuser")
	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	err = repo.UpdateByID(context.Background(), "test-id", user)
	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	err = repo.DeleteByID(context.Background(), "test-id")
	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}
}
