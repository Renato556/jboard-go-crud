package models

import (
	"time"
)

type Job struct {
	ID                      string            `json:"id" bson:"_id" validate:"required"`
	Title                   string            `json:"title" bson:"title" validate:"required"`
	UpdatedAt               string            `json:"updatedAt" bson:"updatedAt"`
	EmploymentType          string            `json:"employmentType" bson:"employmentType"`
	PublishedDate           string            `json:"publishedDate" bson:"publishedDate"`
	ApplicationDeadline     string            `json:"applicationDeadline" bson:"applicationDeadline"`
	CompensationTierSummary string            `json:"compensationTierSummary" bson:"compensationTierSummary"`
	WorkplaceType           string            `json:"workplaceType" bson:"workplaceType"`
	OfficeLocation          string            `json:"officeLocation" bson:"officeLocation"`
	IsBrazilianFriendly     BrazilianFriendly `json:"isBrazilianFriendly" bson:"isBrazilianFriendly"`
	Company                 string            `json:"company" bson:"company" validate:"required"`
	Url                     string            `json:"url" bson:"url" validate:"required"`
	SeniorityLevel          string            `json:"seniorityLevel" bson:"seniorityLevel" validate:"required"`
	Field                   string            `json:"field" bson:"field" validate:"required"`
	ExpiresAt               time.Time         `json:"expiresAt" bson:"expiresAt"`
}
