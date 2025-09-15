package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskCreateEditRequest struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" form:"title" validate:"required,max=100"`
	ProjectID uuid.UUID `json:"project_id,omitempty" form:"project_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type TaskReponseGet struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	ProjectID uuid.UUID `db:"project_id" json:"project_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
