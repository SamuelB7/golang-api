package models

import (
	"time"

	"github.com/google/uuid"
)

type Posts struct {
	ID        uuid.UUID `json:"id" validate:"-"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	UserID    string    `json:"user_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
