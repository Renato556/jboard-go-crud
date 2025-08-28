package repositories

import (
	"context"
	"errors"
	"jboard-go-crud/src/models"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var validate = validator.New()

type JobRepository interface {
	Create(ctx context.Context, job models.Job) error
	FindAll(ctx context.Context) ([]models.Job, error)
	FindByURL(ctx context.Context, url string) (models.Job, bool, error)
	UpdateByURL(ctx context.Context, url string, job models.Job) error
	DeleteByURL(ctx context.Context, url string, job models.Job) error
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
	_ = repo.ensureIndexes(context.Background())
	return repo
}

func (m *mongoJobRepository) ensureIndexes(ctx context.Context) error {
	coll := m.client.Database(m.database).Collection(m.collection)
	model := mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_url"),
	}
	_, err := coll.Indexes().CreateOne(ctx, model)
	return err
}

func (m *mongoJobRepository) FindAll(ctx context.Context) ([]models.Job, error) {
	coll := m.client.Database(m.database).Collection(m.collection)

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Job
	for cur.Next(ctx) {
		var j models.Job
		if err := cur.Decode(&j); err != nil {
			return nil, err
		}
		results = append(results, j)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (m *mongoJobRepository) Create(ctx context.Context, job models.Job) error {
	if err := validate.Struct(job); err != nil {
		return err
	}
	coll := m.client.Database(m.database).Collection(m.collection)
	_, err := coll.InsertOne(ctx, job)
	return err
}

func (m *mongoJobRepository) FindByURL(ctx context.Context, url string) (models.Job, bool, error) {
	coll := m.client.Database(m.database).Collection(m.collection)
	var result models.Job
	err := coll.FindOne(ctx, bson.M{"url": url}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Job{}, false, nil
		}
		return models.Job{}, false, err
	}
	return result, true, nil
}

func (m *mongoJobRepository) UpdateByURL(ctx context.Context, url string, job models.Job) error {
	if err := validate.Struct(job); err != nil {
		return err
	}
	coll := m.client.Database(m.database).Collection(m.collection)
	_, err := coll.ReplaceOne(ctx, bson.M{"url": url}, job)
	return err
}

func (m *mongoJobRepository) DeleteByURL(ctx context.Context, url string, job models.Job) error {
	if err := validate.Struct(job); err != nil {
		return err
	}
	coll := m.client.Database(m.database).Collection(m.collection)
	_, err := coll.DeleteOne(ctx, bson.M{"url": url})
	return err
}
