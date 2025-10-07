package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"jboard-go-crud/internal/models/enums"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	Role     enums.RoleEnum     `json:"role" bson:"role"`
}
