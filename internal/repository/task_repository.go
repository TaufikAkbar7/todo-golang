package repository

import (
	"context"
	"golang-todo/internal/entity"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TaskRepository struct {
	DB  *sqlx.DB
	Log *logrus.Logger
}

func NewTaskRepository(db *sqlx.DB, log *logrus.Logger) *TaskRepository {
	return &TaskRepository{
		DB:  db,
		Log: log,
	}
}

func (c *TaskRepository) GetAll(ctx context.Context) (*[]entity.Task, error) {
	query := "SELECT id, project_id, title, created_at FROM tasks"
	tasks := new([]entity.Task)

	err := c.DB.SelectContext(ctx, tasks, query)

	return tasks, err
}

func (c *TaskRepository) GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*entity.Task, error) {
	query := "SELECT id, project_id, title, created_at FROM tasks WHERE id = $1"
	task := new(entity.Task)

	err := tx.GetContext(ctx, task, query, id)
	return task, err
}

func (c *TaskRepository) Create(ctx context.Context, tx *sqlx.Tx, task *entity.Task) error {
	query := "INSERT INTO tasks (id, project_id, title, created_at, updated_at, due_date, status, priority) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := tx.ExecContext(ctx, query, task.ID, task.ProjectID, task.Title, task.CreatedAt, task.UpdatedAt, task.DueDate, task.Status, task.Priority)
	if err != nil {
		return err
	}
	return nil
}

func (c *TaskRepository) Update(ctx context.Context, tx *sqlx.Tx, entity *entity.Task) error {
	query := "UPDATE tasks SET title = $1, project_id = $2, updated_at = $3, due_date = $5, status = $6, priority = $7 WHERE id = $4"
	_, err := tx.ExecContext(ctx, query, entity.Title, entity.ProjectID, entity.UpdatedAt, entity.ID, entity.DueDate, entity.Status, entity.Priority)
	if err != nil {
		return err
	}
	return nil
}

func (c *TaskRepository) Delete(ctx context.Context, tx *sqlx.Tx, id string) error {
	query := "DELETE FROM tasks WHERE id = $1"

	_, err := tx.ExecContext(ctx, query, id)
	return err
}

func (c *TaskRepository) AssignTag(ctx context.Context, tx *sqlx.Tx, entity *entity.TaskTag) error {
	query := "INSERT INTO task_tags (id, task_id, tag_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"

	_, err := tx.ExecContext(ctx, query, entity.ID, entity.TaskID, entity.TagID, entity.CreatedAt, entity.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (c *TaskRepository) UnassignTag(ctx context.Context, tx *sqlx.Tx, id string) error {
	query := "DELETE FROM task_tags WHERE id = $1"

	_, err := tx.ExecContext(ctx, query, id)
	return err
}

func (c *TaskRepository) GetByIDTag(ctx context.Context, tx *sqlx.Tx, id string) (*entity.TaskTag, error) {
	query := "SELECT * FROM task_tags WHERE id = $1"
	task := new(entity.TaskTag)

	err := tx.GetContext(ctx, task, query, id)
	return task, err
}
