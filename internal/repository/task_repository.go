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
	query := "INSERT INTO tasks (id, project_id, title, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.ExecContext(ctx, query, task.ID, task.ProjectID, task.Title, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (c *TaskRepository) Update(ctx context.Context, tx *sqlx.Tx, entity *entity.Task) error {
	query := "UPDATE tasks SET title = $1, project_id = $2, updated_at = $3 WHERE id = $4"
	_, err := tx.ExecContext(ctx, query, entity.Title, entity.ProjectID, entity.UpdatedAt, entity.ID)
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
