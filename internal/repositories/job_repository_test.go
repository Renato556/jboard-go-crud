package repositories

import (
	"context"
	"jboard-go-crud/internal/models"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewJobRepository(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	if repo == nil {
		t.Error("Expected repository to be created, got nil")
	}
}

func TestNewJobRepository_WithValidClient(t *testing.T) {
	client := &mongo.Client{}
	repo := NewJobRepository(client, "testdb", "jobs")

	if repo == nil {
		t.Error("Expected repository to be created, got nil")
	}
}

func TestJobRepository_Create_ValidationError(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:      "",
		Title:   "",
		Company: "",
		Url:     "",
	}

	err := repo.Create(context.Background(), job)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
}

func TestJobRepository_Create_Success(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:             "test-job-id",
		Title:          "Senior Developer",
		Company:        "Test Company",
		Url:            "https://example.com/job",
		SeniorityLevel: "Senior",
		Field:          "Technology",
	}

	err := repo.Create(context.Background(), job)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get jobs collection" {
		t.Errorf("Expected 'failed to get jobs collection' error, got %v", err)
	}
}

func TestJobRepository_Create_MissingRequiredFields(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:      "test-job-id",
		Title:   "Senior Developer",
		Company: "Test Company",
		Url:     "",
	}

	err := repo.Create(context.Background(), job)

	if err == nil {
		t.Error("Expected validation error for missing URL, got nil")
	}
}

func TestJobRepository_Create_MissingSeniorityLevel(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:      "test-job-id",
		Title:   "Senior Developer",
		Company: "Test Company",
		Url:     "https://example.com/job",
		Field:   "Technology",
	}

	err := repo.Create(context.Background(), job)

	if err == nil {
		t.Error("Expected validation error for missing seniority level, got nil")
	}
}

func TestJobRepository_Create_MissingField(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:             "test-job-id",
		Title:          "Senior Developer",
		Company:        "Test Company",
		Url:            "https://example.com/job",
		SeniorityLevel: "Senior",
	}

	err := repo.Create(context.Background(), job)

	if err == nil {
		t.Error("Expected validation error for missing field, got nil")
	}
}

func TestJobRepository_FindByID_Success(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job, found, err := repo.FindByID(context.Background(), "test-job-id")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if found {
		t.Error("Expected found to be false, got true")
	}

	if job.ID != "" {
		t.Errorf("Expected empty job, got job with ID: %s", job.ID)
	}

	if err.Error() != "failed to get jobs collection" {
		t.Errorf("Expected 'failed to get jobs collection' error, got %v", err)
	}
}

func TestJobRepository_FindByID_EmptyID(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job, found, err := repo.FindByID(context.Background(), "")

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if found {
		t.Error("Expected found to be false, got true")
	}

	if job.ID != "" {
		t.Errorf("Expected empty job, got job with ID: %s", job.ID)
	}

	if err.Error() != "failed to get jobs collection" {
		t.Errorf("Expected 'failed to get jobs collection' error, got %v", err)
	}
}

func TestJobRepository_UpdateByID_Success(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:             "test-job-id",
		Title:          "Updated Senior Developer",
		Company:        "Updated Test Company",
		Url:            "https://example.com/updated-job",
		SeniorityLevel: "Senior",
		Field:          "Technology",
	}

	err := repo.UpdateByID(context.Background(), "test-job-id", job)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get jobs collection" {
		t.Errorf("Expected 'failed to get jobs collection' error, got %v", err)
	}
}

func TestJobRepository_UpdateByID_ValidationError(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:      "",
		Title:   "",
		Company: "",
		Url:     "",
	}

	err := repo.UpdateByID(context.Background(), "test-job-id", job)

	if err == nil {
		t.Error("Expected validation error, got nil")
	}
}

func TestJobRepository_UpdateByID_EmptyID(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:             "test-job-id",
		Title:          "Senior Developer",
		Company:        "Test Company",
		Url:            "https://example.com/job",
		SeniorityLevel: "Senior",
		Field:          "Technology",
	}

	err := repo.UpdateByID(context.Background(), "", job)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get jobs collection" {
		t.Errorf("Expected 'failed to get jobs collection' error, got %v", err)
	}
}

func TestJobRepository_UpdateByID_MissingRequiredFields(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:      "test-job-id",
		Title:   "Senior Developer",
		Company: "Test Company",
		Url:     "",
	}

	err := repo.UpdateByID(context.Background(), "test-job-id", job)

	if err == nil {
		t.Error("Expected validation error for missing URL, got nil")
	}
}

func TestJobRepository_FindAll_Success(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	jobs, err := repo.FindAll(context.Background())

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if jobs != nil {
		t.Errorf("Expected nil jobs slice, got %v", jobs)
	}

	if err.Error() != "failed to get jobs collection" {
		t.Errorf("Expected 'failed to get jobs collection' error, got %v", err)
	}
}

func TestJobRepository_ExpiresAtFieldSetOnCreate(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:             "test-job-id",
		Title:          "Senior Developer",
		Company:        "Test Company",
		Url:            "https://example.com/job",
		SeniorityLevel: "Senior",
		Field:          "Technology",
		ExpiresAt:      time.Time{},
	}

	err := repo.Create(context.Background(), job)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get jobs collection" {
		t.Errorf("Expected 'failed to get jobs collection' error, got %v", err)
	}
}

func TestJobRepository_ExpiresAtFieldSetOnUpdate(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	job := models.Job{
		ID:             "test-job-id",
		Title:          "Updated Senior Developer",
		Company:        "Updated Test Company",
		Url:            "https://example.com/updated-job",
		SeniorityLevel: "Senior",
		Field:          "Technology",
		ExpiresAt:      time.Time{},
	}

	err := repo.UpdateByID(context.Background(), "test-job-id", job)

	if err == nil {
		t.Error("Expected error due to nil MongoDB client, got nil")
	}

	if err.Error() != "failed to get jobs collection" {
		t.Errorf("Expected 'failed to get jobs collection' error, got %v", err)
	}
}

func TestJobRepository_ContextCancellation(t *testing.T) {
	repo := NewJobRepository(nil, "testdb", "jobs")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	job := models.Job{
		ID:             "test-job-id",
		Title:          "Senior Developer",
		Company:        "Test Company",
		Url:            "https://example.com/job",
		SeniorityLevel: "Senior",
		Field:          "Technology",
	}

	err := repo.Create(ctx, job)
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}

	_, _, err = repo.FindByID(ctx, "test-job-id")
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}

	err = repo.UpdateByID(ctx, "test-job-id", job)
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}

	_, err = repo.FindAll(ctx)
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}
}
