package model

import (
	"golang-todo/internal/entity"
	"time"

	"github.com/google/uuid"
)

type TaskCreateEditRequest struct {
	ID        string              `json:"id"`
	Title     string              `json:"title" form:"title" validate:"required,max=100"`
	ProjectID uuid.UUID           `json:"project_id,omitempty" form:"project_id"`
	CreatedAt *time.Time          `json:"created_at,omitempty"`
	UpdatedAt *time.Time          `json:"updated_at,omitempty"`
	DueDate   *string             `json:"due_date,omitempty" form:"due_date"`
	Status    entity.TaskStatus   `json:"status" form:"status" validate:"required,oneof=todo in_progress done"`
	Priority  entity.TaskPriority `json:"priority" form:"priority" validate:"required,oneof=low medium high"`
}

type TaskReponseGet struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	ProjectID uuid.UUID `db:"project_id" json:"project_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
