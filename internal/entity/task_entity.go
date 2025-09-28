package entity

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string
type TaskPriority string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"

	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	ProjectID uuid.UUID    `db:"project_id" json:"project_id"`
	Title     string       `db:"title" json:"title"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
	DueDate   *time.Time   `db:"due_date" json:"due_date"`
	Status    TaskStatus   `db:"status" json:"status"`
	Priority  TaskPriority `db:"priority" json:"priority"`
}
