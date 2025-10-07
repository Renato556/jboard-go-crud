package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Skill struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Skills   []string           `json:"skills" bson:"skills"`
}

type SkillRequest struct {
	Username string `json:"username"`
	Skill    string `json:"skill"`
}
