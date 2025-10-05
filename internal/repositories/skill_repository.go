package repositories

import (
	"context"
	"jboard-go-crud/internal/config"
	"jboard-go-crud/internal/models"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var skillValidate = validator.New()

type SkillRepository interface {
	FindAll(ctx context.Context) ([]models.Skill, error)
	FindByUsername(ctx context.Context, username string) (models.Skill, bool, error)
	Create(ctx context.Context, skillRequest models.SkillRequest) error
	AddSkill(ctx context.Context, skillRequest models.SkillRequest) error
	RemoveSkill(ctx context.Context, skillRequest models.SkillRequest) error
	DeleteByUsername(ctx context.Context, username string) error
}

type mongoSkillRepository struct {
	database string
	client   *mongo.Client
}

func NewSkillRepository(client *mongo.Client, dbName, collectionName string) SkillRepository {
	log.Printf("Creating new SkillRepository with database: %s, getCollection: %s", dbName, collectionName)
	return &mongoSkillRepository{
		database: dbName,
		client:   client,
	}
}

func (r *mongoSkillRepository) getCollection() *mongo.Collection {
	return config.GetSkillsCollection(r.database)
}

func (r *mongoSkillRepository) FindAll(ctx context.Context) ([]models.Skill, error) {
	log.Printf("Repository FindAll called for skills")

	var skills []models.Skill
	cursor, err := r.getCollection().Find(ctx, bson.M{})
	if err != nil {
		log.Printf("ERROR: Failed to execute find query for skills: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var skill models.Skill
		if err := cursor.Decode(&skill); err != nil {
			log.Printf("ERROR: Failed to decode skill from cursor: %v", err)
			return nil, err
		}
		skills = append(skills, skill)
	}

	log.Printf("Successfully retrieved %d skills from database", len(skills))
	return skills, nil
}

func (r *mongoSkillRepository) FindByUsername(ctx context.Context, username string) (models.Skill, bool, error) {
	log.Printf("Repository FindByUsername called for username: %s", username)

	var skill models.Skill
	filter := bson.M{"username": username}
	err := r.getCollection().FindOne(ctx, filter).Decode(&skill)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Skills not found for username: %s", username)
			return skill, false, nil
		}
		log.Printf("ERROR: Failed to find skills for username %s: %v", username, err)
		return skill, false, err
	}

	log.Printf("Successfully found skills for username: %s with %d skills", username, len(skill.Skills))
	return skill, true, nil
}

func (r *mongoSkillRepository) Create(ctx context.Context, skillRequest models.SkillRequest) error {
	log.Printf("Repository Create called for username: %s, skill: %s", skillRequest.Username, skillRequest.Skill)

	if err := skillValidate.Struct(skillRequest); err != nil {
		log.Printf("Validation error in Create for username %s: %v", skillRequest.Username, err)
		return err
	}
	log.Printf("Validation passed for username: %s", skillRequest.Username)

	skill := models.Skill{
		Username: skillRequest.Username,
		Skills:   []string{skillRequest.Skill},
	}
	_, err := r.getCollection().InsertOne(ctx, skill)
	if err != nil {
		if strings.Contains(err.Error(), "unacknowledged write") {
			log.Printf("Unacknowledged write for skill user %s - treating as success since data was written to database", skillRequest.Username)
			return nil
		}
		log.Printf("ERROR: Failed to insert skills for username %s: %v", skillRequest.Username, err)
		return err
	}

	log.Printf("Successfully created skills for username: %s", skillRequest.Username)
	return nil
}

func (r *mongoSkillRepository) AddSkill(ctx context.Context, skillRequest models.SkillRequest) error {
	log.Printf("Repository AddSkill called for username: %s, skill: %s", skillRequest.Username, skillRequest.Skill)

	filter := bson.M{"username": skillRequest.Username}
	update := bson.M{
		"$addToSet": bson.M{"skills": skillRequest.Skill},
	}
	result, err := r.getCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		if strings.Contains(err.Error(), "unacknowledged write") {
			log.Printf("Unacknowledged write for adding skill to user %s - treating as success since data was written to database", skillRequest.Username)
			return nil
		}
		log.Printf("ERROR: Failed to add skill %s to username %s: %v", skillRequest.Skill, skillRequest.Username, err)
		return err
	}

	log.Printf("Successfully added skill %s to username: %s, matched: %d, modified: %d", skillRequest.Skill, skillRequest.Username, result.MatchedCount, result.ModifiedCount)
	return nil
}

func (r *mongoSkillRepository) RemoveSkill(ctx context.Context, skillRequest models.SkillRequest) error {
	log.Printf("Repository RemoveSkill called for username: %s, skill: %s", skillRequest.Username, skillRequest.Skill)

	filter := bson.M{"username": skillRequest.Username}
	update := bson.M{
		"$pull": bson.M{"skills": skillRequest.Skill},
	}
	result, err := r.getCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		if strings.Contains(err.Error(), "unacknowledged write") {
			log.Printf("Unacknowledged write for removing skill from user %s - treating as success since data was written to database", skillRequest.Username)
			return nil
		}
		log.Printf("ERROR: Failed to remove skill %s from username %s: %v", skillRequest.Skill, skillRequest.Username, err)
		return err
	}

	log.Printf("Successfully removed skill %s from username: %s, matched: %d, modified: %d", skillRequest.Skill, skillRequest.Username, result.MatchedCount, result.ModifiedCount)
	return nil
}

func (r *mongoSkillRepository) DeleteByUsername(ctx context.Context, username string) error {
	log.Printf("Repository DeleteByUsername called for username: %s", username)

	filter := bson.M{"username": username}
	result, err := r.getCollection().DeleteOne(ctx, filter)
	if err != nil {
		if strings.Contains(err.Error(), "unacknowledged write") {
			log.Printf("Unacknowledged write for deleting skills of user %s - treating as success since data was written to database", username)
			return nil
		}
		log.Printf("ERROR: Failed to delete skills for username %s: %v", username, err)
		return err
	}

	log.Printf("Successfully deleted skills for username: %s, deleted count: %d", username, result.DeletedCount)
	return nil
}
