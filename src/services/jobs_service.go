package services

import (
	"context"
	"errors"
	"reflect"

	"jboard-go-crud/src/models"
	"jboard-go-crud/src/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrJobNotFound = errors.New("job não encontrado")

type UpsertOutcome int

const (
	OutcomeUpdated   UpsertOutcome = 200
	OutcomeCreated   UpsertOutcome = 201
	OutcomeUnchanged UpsertOutcome = 409
)

type JobService interface {
	CreateOrUpdate(ctx context.Context, job models.Job) (UpsertOutcome, error)
	UpdateOnlyIfURLExists(ctx context.Context, job models.Job) error
	DeleteOnlyIfURLExists(ctx context.Context, job models.Job) error
}

type jobService struct {
	repo repositories.JobRepository
}

func NewJobService(r repositories.JobRepository) JobService {
	return &jobService{repo: r}
}

func (s *jobService) CreateOrUpdate(ctx context.Context, job models.Job) (UpsertOutcome, error) {
	existing, found, err := s.repo.FindByURL(ctx, job.Url)
	if err != nil {
		return 0, err
	}
	if !found {
		if err := s.repo.Create(ctx, job); err != nil {
			return 0, err
		}
		return OutcomeCreated, nil
	}

	if equalIgnoringID(existing, job) {
		return OutcomeUnchanged, nil
	}

	if err := s.repo.UpdateByURL(ctx, job.Url, job); err != nil {
		return 0, err
	}
	return OutcomeUpdated, nil
}

func (s *jobService) UpdateOnlyIfURLExists(ctx context.Context, job models.Job) error {
	_, found, err := s.repo.FindByURL(ctx, job.Url)
	if err != nil {
		return err
	}
	if !found {
		return ErrJobNotFound
	}
	return s.repo.UpdateByURL(ctx, job.Url, job)
}

func (s *jobService) DeleteOnlyIfURLExists(ctx context.Context, job models.Job) error {
	_, found, err := s.repo.FindByURL(ctx, job.Url)
	if err != nil {
		return err
	}
	if !found {
		return ErrJobNotFound
	}
	return s.repo.DeleteByURL(ctx, job.Url, job)
}

func equalIgnoringID(a, b models.Job) bool {
	a2 := a
	b2 := b
	a2.ID = bson.ObjectID{}
	b2.ID = bson.ObjectID{}
	return reflect.DeepEqual(a2, b2)
}
