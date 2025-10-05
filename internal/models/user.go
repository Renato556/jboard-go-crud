package models

import "jboard-go-crud/internal/models/enums"

type User struct {
	ID       string         `json:"id" bson:"_id,omitempty"`
	Username string         `json:"username" bson:"username"`
	Password string         `json:"password" bson:"password"`
	Role     enums.RoleEnum `json:"role" bson:"role"`
}
