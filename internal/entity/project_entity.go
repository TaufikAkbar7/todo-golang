package entity

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"project_name"`
	Description string    `db:"description" json:"description"`
	OwnerID     uuid.UUID `db:"owner_id" json:"owner_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
