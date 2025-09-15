package entity

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ProjectID uuid.UUID `db:"project_id" json:"project_id"`
	Title     string    `db:"title" json:"title"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
