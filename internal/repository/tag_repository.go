package repository

import (
	"context"
	"golang-todo/internal/entity"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TagRepository struct {
	DB  *sqlx.DB
	Log *logrus.Logger
}

func NewTagRepository(db *sqlx.DB, log *logrus.Logger) *TagRepository {
	return &TagRepository{
		DB:  db,
		Log: log,
	}
}

func (c *TagRepository) GetAll(ctx context.Context) (*[]entity.Tag, error) {
	query := "SELECT id, name, color, created_at FROM tags"
	tags := new([]entity.Tag)

	err := c.DB.SelectContext(ctx, tags, query)

	return tags, err
}

func (c *TagRepository) GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*entity.Tag, error) {
	query := "SELECT id, name, color, created_at FROM tags WHERE id = $1"
	task := new(entity.Tag)

	err := tx.GetContext(ctx, task, query, id)
	return task, err
}

func (c *TagRepository) Create(ctx context.Context, tx *sqlx.Tx, task *entity.Tag) error {
	query := "INSERT INTO tags (id, name, color, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.ExecContext(ctx, query, task.ID, task.Name, task.Color, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (c *TagRepository) Update(ctx context.Context, tx *sqlx.Tx, entity *entity.Tag) error {
	query := "UPDATE tags SET name = $1, color = $2, updated_at = $3 WHERE id = $4"
	_, err := tx.ExecContext(ctx, query, entity.Name, entity.Color, entity.UpdatedAt, entity.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *TagRepository) Delete(ctx context.Context, tx *sqlx.Tx, id string) error {
	query := "DELETE FROM tags WHERE id = $1"

	_, err := tx.ExecContext(ctx, query, id)
	return err
}
