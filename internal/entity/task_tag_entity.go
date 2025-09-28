package entity

import (
	"time"

	"github.com/google/uuid"
)

type TaskTag struct {
	ID        uuid.UUID `db:"id" json:"id"`
	TaskID    uuid.UUID `db:"task_id" json:"task_id"`
	TagID     uuid.UUID `db:"tag_id" json:"tag_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
