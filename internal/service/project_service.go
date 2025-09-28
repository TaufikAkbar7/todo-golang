package service

import (
	"context"
	"database/sql"
	"errors"
	"golang-todo/internal/entity"
	"golang-todo/internal/helper"
	"golang-todo/internal/model"
	"golang-todo/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ProjectService struct {
	DB                *sqlx.DB
	Repo              *repository.ProjectRepository
	Log               *logrus.Logger
	RepoProjectMember *repository.ProjectMemberRepository
}

func NewProjectService(db *sqlx.DB, repo *repository.ProjectRepository, log *logrus.Logger, repoProjectMember *repository.ProjectMemberRepository) *ProjectService {
	return &ProjectService{
		DB:                db,
		Repo:              repo,
		Log:               log,
		RepoProjectMember: repoProjectMember,
	}
}

func (c *ProjectService) GetAll(ctx context.Context) (*[]model.ProjectReponseGet, error) {
	projects, err := c.Repo.GetAll(ctx)
	if err != nil {
		c.Log.Errorf("Failed to get all data projects %v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Info("Success get all project")
	return projects, nil
}

func (c *ProjectService) GetByID(ctx context.Context, id string) (*entity.Project, error) {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	project, err := c.Repo.GetByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			c.Log.Info("Data not found")
			return nil, fiber.ErrNotFound
		}
		c.Log.Errorf("Failed to get data project %v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Info("Success get project by id")
	return project, nil
}

func (c *ProjectService) Create(ctx context.Context, req *model.ProjectCreateEditRequest, ownerID string) error {
	newID, _ := uuid.NewV7()
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	parsedUUID, err := uuid.Parse(ownerID)
	if err != nil {
		c.Log.Errorf("Failed to parse UUID string: %v", err)
	}
	dateNow := helper.GetDateNow()
	payload := &entity.Project{
		ID:          newID,
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     parsedUUID,
		CreatedAt:   dateNow,
		UpdatedAt:   dateNow,
	}

	if err := c.Repo.Create(ctx, tx, payload); err != nil {
		c.Log.Errorf("Failed create project %v", err)
		return fiber.ErrInternalServerError
	}

	idProjectMember, _ := uuid.NewV7()
	payloadProjectMember := &entity.ProjectMember{
		ID:        idProjectMember,
		ProjectID: newID,
		UserID:    parsedUUID,
		RoleID:    1, // owner
	}
	if err := c.RepoProjectMember.Create(ctx, tx, payloadProjectMember); err != nil {
		c.Log.Errorf("Failed create project member %v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success created project")
	return nil
}

func (c *ProjectService) Update(ctx context.Context, req *model.ProjectCreateEditRequest, id string) error {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	project, err := c.Repo.GetByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			c.Log.Info("Data not found")
			return fiber.ErrNotFound
		}
		c.Log.Errorf("Failed to get data project %v", err)
		return fiber.ErrInternalServerError
	}

	dateNow := helper.GetDateNow()
	payload := &entity.Project{
		ID:          project.ID,
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     project.OwnerID,
		UpdatedAt:   dateNow,
	}

	if err := c.Repo.Update(ctx, tx, payload); err != nil {
		c.Log.Errorf("Failed create project %v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success updated project")
	return nil
}

func (c *ProjectService) Delete(ctx context.Context, id string) error {
	tx, err := c.DB.BeginTxx(ctx, nil)
	if err != nil {
		c.Log.Errorf("Failed start transaction db %v", err)
		return err
	}

	if _, err := c.Repo.GetByID(ctx, tx, id); err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			c.Log.Info("Data not found")
			return fiber.ErrNotFound
		}
		c.Log.Errorf("Failed to get data project %v", err)
		return fiber.ErrInternalServerError
	}

	// check if project still remaining tasks
	totalTasks, err := c.Repo.CountTaskByProjectID(ctx, tx, id)
	if err != nil {
		c.Log.Errorf("Failed to count task by project %v", err)
	}
	if totalTasks != nil && *totalTasks > 0 {
		c.Log.Infof("The project cannot be deleted because there are still remaining tasks: %v", *totalTasks)
		return fiber.NewError(fiber.StatusInternalServerError, "The project cannot be deleted because there are still remaining tasks")
	}

	if err := c.Repo.Delete(ctx, tx, id); err != nil {
		c.Log.Errorf("Failed to deleted user %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success deleted project")
	return nil
}
