package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectCreateEditRequest struct {
	ID          string     `json:"id"`
	Name        string     `json:"project_name" form:"project_name" validate:"required,max=100"`
	Description string     `json:"description" form:"description" validate:"required,max=100"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type ProjectReponseGet struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"project_name"`
	Description string    `db:"description" json:"description"`
	OwnerID     uuid.UUID `db:"owner_id" json:"owner_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
