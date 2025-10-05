package services

import (
	"context"
	"jboard-go-crud/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSkillRepository struct {
	mock.Mock
}

func (m *MockSkillRepository) FindAll(ctx context.Context) ([]models.Skill, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Skill), args.Error(1)
}

func (m *MockSkillRepository) FindByUsername(ctx context.Context, username string) (models.Skill, bool, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.Skill), args.Bool(1), args.Error(2)
}

func (m *MockSkillRepository) Create(ctx context.Context, skillRequest models.SkillRequest) error {
	args := m.Called(ctx, skillRequest)
	return args.Error(0)
}

func (m *MockSkillRepository) AddSkill(ctx context.Context, skillRequest models.SkillRequest) error {
	args := m.Called(ctx, skillRequest)
	return args.Error(0)
}

func (m *MockSkillRepository) RemoveSkill(ctx context.Context, skillRequest models.SkillRequest) error {
	args := m.Called(ctx, skillRequest)
	return args.Error(0)
}

func (m *MockSkillRepository) DeleteByUsername(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}

func TestSkillService_AddSkill_NewUser(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	service := NewSkillService(mockRepo)

	skillRequest := models.SkillRequest{
		Username: "testuser",
		Skill:    "java",
	}

	mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(models.Skill{}, false, nil)
	mockRepo.On("Create", mock.Anything, skillRequest).Return(nil)

	err := service.AddSkill(context.Background(), skillRequest)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_AddSkill_ExistingUser(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	service := NewSkillService(mockRepo)

	skillRequest := models.SkillRequest{
		Username: "testuser",
		Skill:    "python",
	}

	existingSkill := models.Skill{
		Username: "testuser",
		Skills:   []string{"java"},
	}

	mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(existingSkill, true, nil)
	mockRepo.On("AddSkill", mock.Anything, skillRequest).Return(nil)

	err := service.AddSkill(context.Background(), skillRequest)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_AddSkill_InvalidRequest(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	service := NewSkillService(mockRepo)

	skillRequest := models.SkillRequest{
		Username: "",
		Skill:    "java",
	}

	err := service.AddSkill(context.Background(), skillRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username and skill are required")
}

func TestSkillService_RemoveSkill_UserNotFound(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	service := NewSkillService(mockRepo)

	skillRequest := models.SkillRequest{
		Username: "nonexistent",
		Skill:    "java",
	}

	mockRepo.On("FindByUsername", mock.Anything, "nonexistent").Return(models.Skill{}, false, nil)

	err := service.RemoveSkill(context.Background(), skillRequest)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}

func TestSkillService_GetAllSkills_Success(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	service := NewSkillService(mockRepo)

	expectedSkill := models.Skill{
		Username: "testuser",
		Skills:   []string{"java", "python"},
	}

	mockRepo.On("FindByUsername", mock.Anything, "testuser").Return(expectedSkill, true, nil)

	result, err := service.GetAllSkills(context.Background(), "testuser")
	assert.NoError(t, err)
	assert.Equal(t, expectedSkill, result)
	mockRepo.AssertExpectations(t)
}

func TestSkillService_GetAllSkills_UserNotFound(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	service := NewSkillService(mockRepo)

	mockRepo.On("FindByUsername", mock.Anything, "nonexistent").Return(models.Skill{}, false, nil)

	_, err := service.GetAllSkills(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}

func TestSkillService_GetAllSkills_EmptyUsername(t *testing.T) {
	mockRepo := new(MockSkillRepository)
	service := NewSkillService(mockRepo)

	_, err := service.GetAllSkills(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username is required")
}
