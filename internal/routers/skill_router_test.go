package routers

import (
	"bytes"
	"context"
	"encoding/json"
	"jboard-go-crud/internal/controllers"
	"jboard-go-crud/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSkillService struct {
	mock.Mock
}

func (m *MockSkillService) GetAllSkills(ctx context.Context, username string) (models.Skill, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.Skill), args.Error(1)
}

func (m *MockSkillService) AddSkill(ctx context.Context, skillRequest models.SkillRequest) error {
	args := m.Called(ctx, skillRequest)
	return args.Error(0)
}

func (m *MockSkillService) RemoveSkill(ctx context.Context, skillRequest models.SkillRequest) error {
	args := m.Called(ctx, skillRequest)
	return args.Error(0)
}

func (m *MockSkillService) DeleteUserSkills(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}

func TestSkillRouter_GetAllSkills(t *testing.T) {
	mockService := new(MockSkillService)
	expectedSkill := models.Skill{
		Username: "testuser",
		Skills:   []string{"java", "python"},
	}
	mockService.On("GetAllSkills", mock.Anything, "testuser").Return(expectedSkill, nil)

	handler := controllers.NewSkillHandler(mockService)
	router := NewSkillsController(handler)

	req, _ := http.NewRequest("GET", "/v1/skills?username=testuser", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	mockService.AssertExpectations(t)
}

func TestSkillRouter_AddSkill(t *testing.T) {
	mockService := new(MockSkillService)
	skillRequest := models.SkillRequest{
		Username: "testuser",
		Skill:    "java",
	}
	mockService.On("AddSkill", mock.Anything, skillRequest).Return(nil)

	handler := controllers.NewSkillHandler(mockService)
	router := NewSkillsController(handler)

	body, _ := json.Marshal(skillRequest)
	req, _ := http.NewRequest("POST", "/v1/skills", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestSkillRouter_RemoveSkill(t *testing.T) {
	mockService := new(MockSkillService)
	skillRequest := models.SkillRequest{
		Username: "testuser",
		Skill:    "java",
	}
	mockService.On("RemoveSkill", mock.Anything, skillRequest).Return(nil)

	handler := controllers.NewSkillHandler(mockService)
	router := NewSkillsController(handler)

	body, _ := json.Marshal(skillRequest)
	req, _ := http.NewRequest("PUT", "/v1/skills", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestSkillRouter_DeleteUserSkills(t *testing.T) {
	mockService := new(MockSkillService)
	mockService.On("DeleteUserSkills", mock.Anything, "testuser").Return(nil)

	handler := controllers.NewSkillHandler(mockService)
	router := NewSkillsController(handler)

	req, _ := http.NewRequest("DELETE", "/v1/skills?username=testuser", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}
