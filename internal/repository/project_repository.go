package repository

import (
	"context"
	"golang-todo/internal/entity"
	"golang-todo/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ProjectRepository struct {
	DB  *sqlx.DB
	Log *logrus.Logger
}

func NewProjectRepository(db *sqlx.DB, log *logrus.Logger) *ProjectRepository {
	return &ProjectRepository{
		DB:  db,
		Log: log,
	}
}

func (c *ProjectRepository) GetAll(ctx context.Context) (*[]model.ProjectReponseGet, error) {
	query := "SELECT id, name, description, owner_id, created_at FROM projects"
	projects := new([]model.ProjectReponseGet)

	err := c.DB.SelectContext(ctx, projects, query)

	return projects, err
}

func (c *ProjectRepository) GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*entity.Project, error) {
	query := "SELECT id, name, description, owner_id, created_at FROM projects WHERE id = $1"
	project := new(entity.Project)

	err := tx.GetContext(ctx, project, query, id)
	return project, err
}

func (c *ProjectRepository) Create(ctx context.Context, tx *sqlx.Tx, project *entity.Project) error {
	query := "INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)"
	c.Log.Info(project)
	_, err := tx.ExecContext(ctx, query, project.ID, project.Name, project.Description, project.OwnerID)
	if err != nil {
		return err
	}
	return nil
}

func (c *ProjectRepository) Update(ctx context.Context, tx *sqlx.Tx, entity *entity.Project) error {
	query := "UPDATE projects SET name = $1, description = $2, owner_id = $3, updated_at = $4 WHERE id = $5"
	c.Log.Debug(&entity)
	_, err := tx.ExecContext(ctx, query, entity.Name, entity.Description, entity.OwnerID, entity.UpdatedAt, entity.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *ProjectRepository) Delete(ctx context.Context, tx *sqlx.Tx, id string) error {
	query := "DELETE FROM projects WHERE id = $1"

	_, err := tx.ExecContext(ctx, query, id)
	return err
}

func (c *ProjectRepository) CountTaskByProjectID(ctx context.Context, tx *sqlx.Tx, id string) (*int, error) {
	query := "SELECT COUNT(t.*) from projects AS p JOIN tasks AS t ON p.id = t.project_id WHERE p.id = $1"
	count := new(int)

	err := tx.GetContext(ctx, count, query, id)
	return count, err
}
