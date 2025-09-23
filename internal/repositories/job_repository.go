package repositories

import (
	"context"
	"errors"
	"jboard-go-crud/internal/models"
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
	client     *mongo.Client
	database   string
	collection string
}

func NewJobRepository(client *mongo.Client, dbName, collectionName string) JobRepository {
	repo := &mongoJobRepository{
		client:     client,
		database:   dbName,
		collection: collectionName,
	}
	if client != nil {
		_ = repo.ensureIndexes(context.Background())
	}
	return repo
}

func (m *mongoJobRepository) ensureIndexes(ctx context.Context) error {
	if m.client == nil {
		return errors.New("MongoDB client is nil")
	}

	coll := m.client.Database(m.database).Collection(m.collection)

	ttlModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err := coll.Indexes().CreateOne(ctx, ttlModel)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoJobRepository) Create(ctx context.Context, job models.Job) error {
	if err := validate.Struct(job); err != nil {
		return err
	}

	job.ExpiresAt = time.Now().Add(12*time.Hour + 1*time.Minute)

	if m.client == nil {
		return errors.New("MongoDB client is nil")
	}

	coll := m.client.Database(m.database).Collection(m.collection)
	_, err := coll.InsertOne(ctx, job)
	return err
}

func (m *mongoJobRepository) FindByID(ctx context.Context, id string) (models.Job, bool, error) {
	if m.client == nil {
		return models.Job{}, false, errors.New("MongoDB client is nil")
	}

	coll := m.client.Database(m.database).Collection(m.collection)
	var result models.Job
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Job{}, false, nil
		}
		return models.Job{}, false, err
	}
	return result, true, nil
}

func (m *mongoJobRepository) UpdateByID(ctx context.Context, id string, job models.Job) error {
	if err := validate.Struct(job); err != nil {
		return err
	}

	job.ExpiresAt = time.Now().Add(12*time.Hour + 1*time.Minute)

	if m.client == nil {
		return errors.New("MongoDB client is nil")
	}

	coll := m.client.Database(m.database).Collection(m.collection)
	_, err := coll.ReplaceOne(ctx, bson.M{"_id": id}, job)
	return err
}

func (m *mongoJobRepository) FindAll(ctx context.Context) ([]models.Job, error) {
	if m.client == nil {
		return nil, errors.New("MongoDB client is nil")
	}

	coll := m.client.Database(m.database).Collection(m.collection)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []models.Job
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
