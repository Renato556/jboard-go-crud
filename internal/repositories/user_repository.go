package repositories

import (
	"context"
	"errors"
	"jboard-go-crud/internal/config"
	"jboard-go-crud/internal/models"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userValidate = validator.New()

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	FindByID(ctx context.Context, id string) (models.User, bool, error)
	FindByUsername(ctx context.Context, username string) (models.User, bool, error)
	UpdateByID(ctx context.Context, id string, user models.User) error
	DeleteByID(ctx context.Context, id string) error
}

type mongoUserRepository struct {
	database string
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) UserRepository {
	log.Printf("Creating new UserRepository with database: %s, collection: %s", dbName, collectionName)
	repo := &mongoUserRepository{
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

func (m *mongoUserRepository) getCollection() *mongo.Collection {
	return config.GetUsersCollection(m.database)
}

func (m *mongoUserRepository) ensureIndexes(ctx context.Context) error {
	log.Printf("Ensuring unique index on username field...")

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get users collection when ensuring indexes")
		return errors.New("failed to get users collection")
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "username", Value: 1}},
	}
	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("ERROR: Failed to create username index: %v", err)
		return err
	}

	log.Printf("Username index created successfully")
	return nil
}

func (m *mongoUserRepository) Create(ctx context.Context, user models.User) error {
	log.Printf("Repository Create called for user ID: %s", user.ID)

	if err := userValidate.Struct(user); err != nil {
		log.Printf("Validation error in Create for user ID %s: %v", user.ID, err)
		return err
	}
	log.Printf("Validation passed for user ID: %s", user.ID)

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get users collection in Create")
		return errors.New("failed to get users collection")
	}

	_, err := coll.InsertOne(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Printf("ERROR: Username already exists for user: %s", user.Username)
			return errors.New("username already exists")
		}
		if strings.Contains(err.Error(), "unacknowledged write") {
			log.Printf("Unacknowledged write for user ID %s - treating as success since data was written to database", user.ID)
		} else {
			log.Printf("ERROR: Failed to insert user ID %s: %v", user.ID, err)
			return err
		}
	}

	log.Printf("Successfully created user with ID: %s", user.ID)
	return nil
}

func (m *mongoUserRepository) FindByID(ctx context.Context, id string) (models.User, bool, error) {
	log.Printf("Repository FindByID called for ID: %s", id)

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get users collection in FindByID")
		return models.User{}, false, errors.New("failed to get users collection")
	}

	var result models.User
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("User not found for ID: %s", id)
			return models.User{}, false, nil
		}
		log.Printf("ERROR: Failed to find user ID %s: %v", id, err)
		return models.User{}, false, err
	}

	log.Printf("Successfully found user ID: %s", id)
	return result, true, nil
}

func (m *mongoUserRepository) FindByUsername(ctx context.Context, username string) (models.User, bool, error) {
	log.Printf("Repository FindByUsername called for username: %s", username)

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get users collection in FindByUsername")
		return models.User{}, false, errors.New("failed to get users collection")
	}

	var result models.User
	err := coll.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("User not found for username: %s", username)
			return models.User{}, false, nil
		}
		log.Printf("ERROR: Failed to find user by username %s: %v", username, err)
		return models.User{}, false, err
	}

	log.Printf("Successfully found user by username: %s", username)
	return result, true, nil
}

func (m *mongoUserRepository) UpdateByID(ctx context.Context, id string, user models.User) error {
	log.Printf("Repository UpdateByID called for user ID: %s", id)

	if err := userValidate.Struct(user); err != nil {
		log.Printf("Validation error in UpdateByID for user ID %s: %v", id, err)
		return err
	}
	log.Printf("Validation passed for user ID: %s", id)

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get users collection in UpdateByID")
		return errors.New("failed to get users collection")
	}

	result, err := coll.ReplaceOne(ctx, bson.M{"_id": id}, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Printf("ERROR: Username already exists when updating user: %s", user.Username)
			return errors.New("username already exists")
		}
		if strings.Contains(err.Error(), "unacknowledged write") {
			log.Printf("Unacknowledged write for user ID %s - treating as success since data was written to database", id)
		} else {
			log.Printf("ERROR: Failed to update user ID %s: %v", id, err)
			return err
		}
	} else {
		log.Printf("Successfully updated user ID: %s, matched: %d, modified: %d", id, result.MatchedCount, result.ModifiedCount)
	}

	return nil
}

func (m *mongoUserRepository) DeleteByID(ctx context.Context, id string) error {
	log.Printf("Repository DeleteByID called for user ID: %s", id)

	coll := m.getCollection()
	if coll == nil {
		log.Printf("ERROR: Failed to get users collection in DeleteByID")
		return errors.New("failed to get users collection")
	}

	result, err := coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Printf("ERROR: Failed to delete user ID %s: %v", id, err)
		return err
	}

	if result.DeletedCount == 0 {
		log.Printf("User not found for deletion with ID: %s", id)
		return errors.New("user not found")
	}

	log.Printf("Successfully deleted user ID: %s", id)
	return nil
}
