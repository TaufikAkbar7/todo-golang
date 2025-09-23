package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMemberCreateEditRequest struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id,omitempty" form:"project_id"`
	RoleID    int       `json:"role_id" form:"role_id" validate:"required,max=100"`
	UserID    uuid.UUID `json:"user_id,omitempty" form:"user_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
