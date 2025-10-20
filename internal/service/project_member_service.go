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

type ProjectMemberService struct {
	DB   *sqlx.DB
	Repo *repository.ProjectMemberRepository
	Log  *logrus.Logger
}

func NewProjectMemberService(db *sqlx.DB, repo *repository.ProjectMemberRepository, log *logrus.Logger) *ProjectMemberService {
	return &ProjectMemberService{
		DB:   db,
		Repo: repo,
		Log:  log,
	}
}

func (c *ProjectMemberService) Create(ctx context.Context, req *model.ProjectMemberCreateEditRequest) error {
	newID, _ := uuid.NewV7()
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	// check if member exist in project
	hasMember, err := c.Repo.FindProjectMember(ctx, tx, req.ProjectID, *req.UserID)
	if err != nil {
		c.Log.Errorf("Failed find member in project %v", err)
		return fiber.ErrInternalServerError
	}
	if hasMember != nil && *hasMember {
		message := "Member already join project"
		c.Log.Info(message)
		return fiber.NewError(400, message)
	}

	dateNow := helper.GetDateNow()
	payload := &entity.ProjectMember{
		ID:        newID,
		ProjectID: req.ProjectID,
		UserID:    *req.UserID,
		RoleID:    req.RoleID,
		CreatedAt: dateNow,
		UpdatedAt: dateNow,
	}
	if err := c.Repo.Create(ctx, tx, payload); err != nil {
		c.Log.Errorf("Failed create project member %v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success invite member to the project")
	return nil
}

func (c *ProjectMemberService) Delete(ctx context.Context, id uuid.UUID, userId uuid.UUID) error {
	tx, err := c.DB.BeginTxx(ctx, nil)
	filters := new(model.ProjectMemberFilter)
	if err != nil {
		c.Log.Errorf("Failed start transaction db %v", err)
		return err
	}

	filters.UserID = &userId
	if _, err := c.Repo.GetByID(ctx, tx, *filters); err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			c.Log.Info("Data not found")
			return fiber.ErrNotFound
		}
		c.Log.Errorf("Failed to get data project member %v", err)
		return fiber.ErrInternalServerError
	}

	if err := c.Repo.Delete(ctx, tx, id, userId); err != nil {
		c.Log.Errorf("Failed to deleted member %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success deleted member")
	return nil
}

func (c *ProjectMemberService) GetByID(ctx context.Context, userId uuid.UUID, projectId uuid.UUID) (*entity.ProjectMember, error) {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()
	filters := new(model.ProjectMemberFilter)

	filters.UserID = &userId
	filters.ProjectID = &projectId
	member, err := c.Repo.GetByID(ctx, tx, *filters)
	if err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			c.Log.Info("Data not found")
			return nil, fiber.ErrNotFound
		}
		c.Log.Errorf("Failed to get data project member %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Info("Success get member")
	return member, nil
}
