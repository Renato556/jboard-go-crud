package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func TestSkillHandler_GetAllSkills_Success(t *testing.T) {
	mockService := new(MockSkillService)
	handler := NewSkillHandler(mockService)

	expectedSkill := models.Skill{
		ID:       "1",
		Username: "testuser",
		Skills:   []string{"java", "python"},
	}

	mockService.On("GetAllSkills", mock.Anything, "testuser").Return(expectedSkill, nil)

	req, _ := http.NewRequest("GET", "/v1/skills?username=testuser", nil)
	rr := httptest.NewRecorder()

	handler.GetAllSkills(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var actualSkill models.Skill
	json.Unmarshal(rr.Body.Bytes(), &actualSkill)
	assert.Equal(t, expectedSkill, actualSkill)
	mockService.AssertExpectations(t)
}

func TestSkillHandler_GetAllSkills_MissingUsername(t *testing.T) {
	mockService := new(MockSkillService)
	handler := NewSkillHandler(mockService)

	req, _ := http.NewRequest("GET", "/v1/skills", nil)
	rr := httptest.NewRecorder()

	handler.GetAllSkills(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSkillHandler_GetAllSkills_UserNotFound(t *testing.T) {
	mockService := new(MockSkillService)
	handler := NewSkillHandler(mockService)

	mockService.On("GetAllSkills", mock.Anything, "nonexistent").Return(models.Skill{}, errors.New("user not found"))

	req, _ := http.NewRequest("GET", "/v1/skills?username=nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.GetAllSkills(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockService.AssertExpectations(t)
}

func TestSkillHandler_AddSkill_Success(t *testing.T) {
	mockService := new(MockSkillService)
	handler := NewSkillHandler(mockService)

	skillRequest := models.SkillRequest{
		Username: "testuser",
		Skill:    "java",
	}

	mockService.On("AddSkill", mock.Anything, skillRequest).Return(nil)

	body, _ := json.Marshal(skillRequest)
	req, _ := http.NewRequest("POST", "/v1/skills", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.AddSkill(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestSkillHandler_AddSkill_InvalidJSON(t *testing.T) {
	mockService := new(MockSkillService)
	handler := NewSkillHandler(mockService)

	req, _ := http.NewRequest("POST", "/v1/skills", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.AddSkill(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSkillHandler_RemoveSkill_Success(t *testing.T) {
	mockService := new(MockSkillService)
	handler := NewSkillHandler(mockService)

	skillRequest := models.SkillRequest{
		Username: "testuser",
		Skill:    "java",
	}

	mockService.On("RemoveSkill", mock.Anything, skillRequest).Return(nil)

	body, _ := json.Marshal(skillRequest)
	req, _ := http.NewRequest("PUT", "/v1/skills", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.RemoveSkill(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestSkillHandler_DeleteUserSkills_Success(t *testing.T) {
	mockService := new(MockSkillService)
	handler := NewSkillHandler(mockService)

	mockService.On("DeleteUserSkills", mock.Anything, "testuser").Return(nil)

	req, _ := http.NewRequest("DELETE", "/v1/skills?username=testuser", nil)
	rr := httptest.NewRecorder()

	handler.DeleteUserSkills(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestSkillHandler_DeleteUserSkills_MissingUsername(t *testing.T) {
	mockService := new(MockSkillService)
	handler := NewSkillHandler(mockService)

	req, _ := http.NewRequest("DELETE", "/v1/skills", nil)
	rr := httptest.NewRecorder()

	handler.DeleteUserSkills(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
