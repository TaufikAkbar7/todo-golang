package repository

import (
	"context"
	"fmt"
	"golang-todo/internal/entity"
	"golang-todo/internal/model"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ProjectMemberRepository struct {
	DB  *sqlx.DB
	Log *logrus.Logger
}

func NewProjectMemberRepository(db *sqlx.DB, log *logrus.Logger) *ProjectMemberRepository {
	return &ProjectMemberRepository{
		DB:  db,
		Log: log,
	}
}

func (c *ProjectMemberRepository) Create(ctx context.Context, tx *sqlx.Tx, project *entity.ProjectMember) error {
	query := "INSERT INTO project_members (id, project_id, role_id, user_id) VALUES ($1, $2, $3, $4)"
	c.Log.Info(project)
	_, err := tx.ExecContext(ctx, query, project.ID, project.ProjectID, project.RoleID, project.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (c *ProjectMemberRepository) Delete(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, userID uuid.UUID) error {
	query := "DELETE FROM project_members WHERE project_id = $1 AND user_id = $2"

	_, err := tx.ExecContext(ctx, query, id, userID)
	return err
}

func (c *ProjectMemberRepository) FindProjectMember(ctx context.Context, tx *sqlx.Tx, projectID uuid.UUID, userID uuid.UUID) (*bool, error) {
	query := `SELECT EXISTS (
    SELECT 1
    FROM project_members
    WHERE user_id = $1
      AND project_id = $2
) AS has_member;
`
	hasMember := new(bool)

	err := tx.GetContext(ctx, hasMember, query, userID, projectID)
	return hasMember, err
}

func (c *ProjectMemberRepository) GetByID(ctx context.Context, tx *sqlx.Tx, filter model.ProjectMemberFilter) (*entity.ProjectMember, error) {
	query := "SELECT * FROM project_members"
	var conditions []string
	var args []any

	// dynamically build the query based on the filter
	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, *filter.UserID)
	}

	if filter.ProjectID != nil {
		conditions = append(conditions, fmt.Sprintf("project_id = $%d", len(args)+1))
		args = append(args, *filter.ProjectID)
	}

	// if any conditions were added, append them to the query
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	project := new(entity.ProjectMember)

	err := tx.GetContext(ctx, project, query, args...)
	return project, err
}
