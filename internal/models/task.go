package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	Priority    int                `bson:"priority" json:"priority"`
	Status      string             `bson:"status" json:"status"`
	DueDate     time.Time          `bson:"due_date,omitempty" json:"due_date"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Priority    int    `json:"priority" validate:"gte=1,lte=5"`
	Status      string `json:"status" validate:"required,oneof=pending completed in_progress"`
	DueDate     string `json:"due_date" validate:"omitempty,rfc3339"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=3"`
	Description *string `json:"description" validate:"omitempty"`
	Category    *string `json:"category" validate:"omitempty"`
	Status      *string `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
	Priority    *int    `json:"priority" validate:"omitempty,gte=1,lte=5"`
	DueDate     *string `json:"due_date" validate:"omitempty,rfc3339"`
}

func (u UpdateTaskRequest) HasUpdates() bool {
	return u.Title != nil ||
		u.Description != nil ||
		u.Status != nil ||
		u.Priority != nil ||
		u.DueDate != nil
}
