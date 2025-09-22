package models

type BrazilianFriendly struct {
	IsFriendly bool   `json:"isFriendly" bson:"isFriendly"`
	Reason     string `json:"reason" bson:"reason"`
}
