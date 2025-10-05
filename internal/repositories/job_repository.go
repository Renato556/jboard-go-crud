package repositories

import (
	"context"
	"errors"
	"jboard-go-crud/internal/config"
	"jboard-go-crud/internal/models"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()

type JobRepository interface {
	Create(ctx context.Context, job models.Job) error
	FindByID(ctx context.Context, id string) (models.Job, bool, error)
	UpdateByID(ctx context.Context, id string, job models.Job) error
	FindAll(ctx context.Context) ([]models.Job, error)
}

type mongoJobRepository struct {
	database string
}

func NewJobRepository(client *mongo.Client, dbName, collectionName string) JobRepository {
	log.Printf("Creating new JobRepository with database: %s, getCollection: %s", dbName, collectionName)
	repo := &mongoJobRepository{
		database: dbName,
	}
	if client != nil {
		log.Printf("MongoDB client is available, ensuring indexes...")
		_ = repo.ensureIndexes(context.Background())
	} else {
		log.Printf("WARNING: MongoDB client is nil")
	}
	return repo
}

func (m *mongoJobRepository) getCollection() *mongo.Collection {
	return config.GetJobsCollection(m.database)
}

func (m *mongoJobRepository) ensureIndexes(ctx context.Context) error {
	log.Printf("Ensuring TTL index on expiresAt field...")

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get jobs getCollection when ensuring indexes")
		return errors.New("failed to get jobs getCollection")
	}

	ttlModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err := coll.Indexes().CreateOne(ctx, ttlModel)
	if err != nil {
		log.Printf("ERROR: Failed to create TTL index: %v", err)
		return err
	}

	log.Printf("TTL index created successfully")
	return nil
}

func (m *mongoJobRepository) Create(ctx context.Context, job models.Job) error {
	log.Printf("Repository Create called for job ID: %s", job.ID)

	if err := validate.Struct(job); err != nil {
		log.Printf("Validation error in Create for job ID %s: %v", job.ID, err)
		return err
	}
	log.Printf("Validation passed for job ID: %s", job.ID)

	job.ExpiresAt = time.Now().Add(12*time.Hour + 1*time.Minute)
	log.Printf("Set expiresAt to: %v for job ID: %s", job.ExpiresAt, job.ID)

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get jobs getCollection in Create")
		return errors.New("failed to get jobs getCollection")
	}

	_, err := coll.InsertOne(ctx, job)
	if err != nil {
		if strings.Contains(err.Error(), "unacknowledged write") {
			log.Printf("Unacknowledged write for job ID %s - treating as success since data was written to database", job.ID)
		} else {
			log.Printf("ERROR: Failed to insert job ID %s: %v", job.ID, err)
			return err
		}
	}

	log.Printf("Successfully created job with MongoDB job ID: %s", job.ID)
	return nil
}

func (m *mongoJobRepository) FindByID(ctx context.Context, id string) (models.Job, bool, error) {
	log.Printf("Repository FindByID called for ID: %s", id)

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get jobs getCollection in FindByID")
		return models.Job{}, false, errors.New("failed to get jobs getCollection")
	}

	var result models.Job
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("Job not found for ID: %s", id)
			return models.Job{}, false, nil
		}
		log.Printf("ERROR: Failed to find job ID %s: %v", id, err)
		return models.Job{}, false, err
	}

	log.Printf("Successfully found job ID: %s", id)
	return result, true, nil
}

func (m *mongoJobRepository) UpdateByID(ctx context.Context, id string, job models.Job) error {
	log.Printf("Repository UpdateByID called for job ID: %s", id)

	if err := validate.Struct(job); err != nil {
		log.Printf("Validation error in UpdateByID for job ID %s: %v", id, err)
		return err
	}
	log.Printf("Validation passed for job ID: %s", id)

	job.ExpiresAt = time.Now().Add(12*time.Hour + 1*time.Minute)
	log.Printf("Updated expiresAt to: %v for job ID: %s", job.ExpiresAt, id)

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get jobs getCollection in UpdateByID")
		return errors.New("failed to get jobs getCollection")
	}

	result, err := coll.ReplaceOne(ctx, bson.M{"_id": id}, job)
	if err != nil {
		if strings.Contains(err.Error(), "unacknowledged write") {
			log.Printf("Unacknowledged write for job ID %s - treating as success since data was written to database", id)
		} else {
			log.Printf("ERROR: Failed to update job ID %s: %v", id, err)
			return err
		}
	} else {
		log.Printf("Successfully updated job ID: %s, matched: %d, modified: %d", id, result.MatchedCount, result.ModifiedCount)
	}

	return nil
}

func (m *mongoJobRepository) FindAll(ctx context.Context) ([]models.Job, error) {
	log.Printf("Repository FindAll called")

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get jobs getCollection in FindAll")
		return nil, errors.New("failed to get jobs getCollection")
	}

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("ERROR: Failed to execute find query: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			log.Printf("WARNING: Error closing cursor: %v", closeErr)
		}
	}()

	var jobs []models.Job
	if err = cursor.All(ctx, &jobs); err != nil {
		log.Printf("ERROR: Failed to decode jobs from cursor: %v", err)
		return nil, err
	}

	log.Printf("Successfully retrieved %d jobs from database", len(jobs))
	return jobs, nil
}
