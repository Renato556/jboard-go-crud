package services

import (
	"context"

	"jboard-go-crud/internal/models"
	"jboard-go-crud/internal/repositories"
)

type UpsertOutcome int

const (
	OutcomeUpdated UpsertOutcome = 200
	OutcomeCreated UpsertOutcome = 201
)

type JobService interface {
	CreateOrUpdate(ctx context.Context, job models.Job) (UpsertOutcome, error)
	FindAll(ctx context.Context) ([]models.Job, error)
}

type jobService struct {
	repo repositories.JobRepository
}

func NewJobService(r repositories.JobRepository) JobService {
	return &jobService{repo: r}
}

func (s *jobService) CreateOrUpdate(ctx context.Context, job models.Job) (UpsertOutcome, error) {
	_, found, err := s.repo.FindByID(ctx, job.ID)
	if err != nil {
		return 0, err
	}

	if found {
		if err := s.repo.UpdateByID(ctx, job.ID, job); err != nil {
			return 0, err
		}
		return OutcomeUpdated, nil
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return 0, err
	}
	return OutcomeCreated, nil
}

func (s *jobService) FindAll(ctx context.Context) ([]models.Job, error) {
	return s.repo.FindAll(ctx)
}
