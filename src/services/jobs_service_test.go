package services

import (
	"context"
	"errors"
	"jboard-go-crud/src/models"
	"testing"
)

type mockJobRepository struct {
	createFunc     func(ctx context.Context, job models.Job) error
	findByIDFunc   func(ctx context.Context, id string) (models.Job, bool, error)
	updateByIDFunc func(ctx context.Context, id string, job models.Job) error
	findAllFunc    func(ctx context.Context) ([]models.Job, error)
}

func (m *mockJobRepository) Create(ctx context.Context, job models.Job) error {
	return m.createFunc(ctx, job)
}

func (m *mockJobRepository) FindByID(ctx context.Context, id string) (models.Job, bool, error) {
	return m.findByIDFunc(ctx, id)
}

func (m *mockJobRepository) UpdateByID(ctx context.Context, id string, job models.Job) error {
	return m.updateByIDFunc(ctx, id, job)
}

func (m *mockJobRepository) FindAll(ctx context.Context) ([]models.Job, error) {
	return m.findAllFunc(ctx)
}

func TestNewJobService(t *testing.T) {
	mockRepo := &mockJobRepository{}
	service := NewJobService(mockRepo)

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}
}

func TestJobService_CreateOrUpdate_Create(t *testing.T) {
	job := models.Job{
		ID:             "test-id",
		Title:          "Test Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
	}

	mockRepo := &mockJobRepository{
		findByIDFunc: func(ctx context.Context, id string) (models.Job, bool, error) {
			return models.Job{}, false, nil
		},
		createFunc: func(ctx context.Context, job models.Job) error {
			return nil
		},
	}

	service := NewJobService(mockRepo)
	ctx := context.Background()

	outcome, err := service.CreateOrUpdate(ctx, job)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if outcome != OutcomeCreated {
		t.Errorf("Expected OutcomeCreated, got %v", outcome)
	}
}

func TestJobService_CreateOrUpdate_Update(t *testing.T) {
	job := models.Job{
		ID:             "test-id",
		Title:          "Updated Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
	}

	existingJob := models.Job{
		ID:             "test-id",
		Title:          "Old Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
	}

	mockRepo := &mockJobRepository{
		findByIDFunc: func(ctx context.Context, id string) (models.Job, bool, error) {
			return existingJob, true, nil
		},
		updateByIDFunc: func(ctx context.Context, id string, job models.Job) error {
			return nil
		},
	}

	service := NewJobService(mockRepo)
	ctx := context.Background()

	outcome, err := service.CreateOrUpdate(ctx, job)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if outcome != OutcomeUpdated {
		t.Errorf("Expected OutcomeUpdated, got %v", outcome)
	}
}

func TestJobService_CreateOrUpdate_FindError(t *testing.T) {
	job := models.Job{
		ID:             "test-id",
		Title:          "Test Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
	}

	mockRepo := &mockJobRepository{
		findByIDFunc: func(ctx context.Context, id string) (models.Job, bool, error) {
			return models.Job{}, false, errors.New("database error")
		},
	}

	service := NewJobService(mockRepo)
	ctx := context.Background()

	_, err := service.CreateOrUpdate(ctx, job)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected 'database error', got %s", err.Error())
	}
}

func TestJobService_CreateOrUpdate_CreateError(t *testing.T) {
	job := models.Job{
		ID:             "test-id",
		Title:          "Test Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
	}

	mockRepo := &mockJobRepository{
		findByIDFunc: func(ctx context.Context, id string) (models.Job, bool, error) {
			return models.Job{}, false, nil
		},
		createFunc: func(ctx context.Context, job models.Job) error {
			return errors.New("create error")
		},
	}

	service := NewJobService(mockRepo)
	ctx := context.Background()

	_, err := service.CreateOrUpdate(ctx, job)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "create error" {
		t.Errorf("Expected 'create error', got %s", err.Error())
	}
}

func TestJobService_CreateOrUpdate_UpdateError(t *testing.T) {
	job := models.Job{
		ID:             "test-id",
		Title:          "Updated Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
	}

	existingJob := models.Job{
		ID:             "test-id",
		Title:          "Old Job",
		Company:        "Test Company",
		Url:            "https://test.com",
		SeniorityLevel: "Senior",
		Field:          "Engineering",
	}

	mockRepo := &mockJobRepository{
		findByIDFunc: func(ctx context.Context, id string) (models.Job, bool, error) {
			return existingJob, true, nil
		},
		updateByIDFunc: func(ctx context.Context, id string, job models.Job) error {
			return errors.New("update error")
		},
	}

	service := NewJobService(mockRepo)
	ctx := context.Background()

	_, err := service.CreateOrUpdate(ctx, job)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "update error" {
		t.Errorf("Expected 'update error', got %s", err.Error())
	}
}

func TestJobService_FindAll_Success(t *testing.T) {
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

	mockRepo := &mockJobRepository{
		findAllFunc: func(ctx context.Context) ([]models.Job, error) {
			return expectedJobs, nil
		},
	}

	service := NewJobService(mockRepo)
	ctx := context.Background()

	jobs, err := service.FindAll(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(jobs))
	}

	if jobs[0].ID != "test-id-1" {
		t.Errorf("Expected first job ID to be 'test-id-1', got %s", jobs[0].ID)
	}

	if jobs[1].ID != "test-id-2" {
		t.Errorf("Expected second job ID to be 'test-id-2', got %s", jobs[1].ID)
	}
}

func TestJobService_FindAll_Error(t *testing.T) {
	mockRepo := &mockJobRepository{
		findAllFunc: func(ctx context.Context) ([]models.Job, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewJobService(mockRepo)
	ctx := context.Background()

	_, err := service.FindAll(ctx)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected 'database error', got %s", err.Error())
	}
}
