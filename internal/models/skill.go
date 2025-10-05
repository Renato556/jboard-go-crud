package models

type Skill struct {
	ID       string   `json:"id" bson:"_id,omitempty"`
	Username string   `json:"username" bson:"username"`
	Skills   []string `json:"skills" bson:"skills"`
}

type SkillRequest struct {
	Username string `json:"username"`
	Skill    string `json:"skill"`
}
