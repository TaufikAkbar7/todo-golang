package model

import (
	"time"
)

type TaskTagCreateEditRequest struct {
	ID        string     `json:"id"`
	TagID     string     `json:"tag_id" form:"tag_id" validate:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
