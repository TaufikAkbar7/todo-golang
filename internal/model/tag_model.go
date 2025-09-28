package model

import (
	"time"

	"github.com/google/uuid"
)

type TagCreateEditRequest struct {
	ID        string     `json:"id"`
	Name      string     `json:"name" form:"name" validate:"required,max=100"`
	Color     string     `json:"color" form:"color" validate:"required,max=7"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type TagReponseGet struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Color     string    `db:"color" json:"color"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
