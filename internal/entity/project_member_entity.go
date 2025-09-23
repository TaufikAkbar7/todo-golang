package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMember struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ProjectID uuid.UUID `db:"project_id" json:"project_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	RoleID    int       `db:"role_id" json:"role_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
