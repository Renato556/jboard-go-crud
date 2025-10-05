package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkillRepository_Interface(t *testing.T) {
	var repo SkillRepository
	mongoRepo := &mongoSkillRepository{
		database: "test",
		client:   nil,
	}
	repo = mongoRepo
	assert.NotNil(t, repo)
}

func TestSkillRepository_Collection(t *testing.T) {
	repo := &mongoSkillRepository{
		database: "test",
		client:   nil,
	}

	assert.NotNil(t, repo)
	assert.Equal(t, "test", repo.database)
}
