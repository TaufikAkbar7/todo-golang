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

type TagService struct {
	DB   *sqlx.DB
	Repo *repository.TagRepository
	Log  *logrus.Logger
}

func NewTagService(db *sqlx.DB, repo *repository.TagRepository, log *logrus.Logger) *TagService {
	return &TagService{
		DB:   db,
		Repo: repo,
		Log:  log,
	}
}

func (c *TagService) GetAll(ctx context.Context) (*[]entity.Tag, error) {
	tags, err := c.Repo.GetAll(ctx)
	if err != nil {
		c.Log.Errorf("Failed to get all data tags %v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Info("Success get all tags")
	return tags, nil
}

func (c *TagService) GetByID(ctx context.Context, id string) (*entity.Tag, error) {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	tag, err := c.Repo.GetByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			c.Log.Info("Data not found")
			return nil, fiber.ErrNotFound
		}
		c.Log.Errorf("Failed to get data tag %v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Info("Success get tag by id")
	return tag, nil
}

func (c *TagService) Create(ctx context.Context, req *model.TagCreateEditRequest) error {
	newID, _ := uuid.NewV7()
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	dateNow := helper.GetDateNow()
	payload := &entity.Tag{
		ID:        newID,
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: dateNow,
		UpdatedAt: dateNow,
	}

	if err := c.Repo.Create(ctx, tx, payload); err != nil {
		c.Log.Errorf("Failed create tag %v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success created tag")
	return nil
}

func (c *TagService) Update(ctx context.Context, req *model.TagCreateEditRequest, id string) error {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	tag, err := c.Repo.GetByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			c.Log.Info("Data not found")
			return fiber.ErrNotFound
		}
		c.Log.Errorf("Failed to get data tag %v", err)
		return fiber.ErrInternalServerError
	}

	dateNow := helper.GetDateNow()
	payload := &entity.Tag{
		ID:        tag.ID,
		Name:      req.Name,
		Color:     req.Color,
		UpdatedAt: dateNow,
	}

	if err := c.Repo.Update(ctx, tx, payload); err != nil {
		c.Log.Errorf("Failed create tag %v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success updated tag")
	return nil
}

func (c *TagService) Delete(ctx context.Context, id string) error {
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
		c.Log.Errorf("Failed to get data tag %v", err)
		return fiber.ErrInternalServerError
	}

	if err := c.Repo.Delete(ctx, tx, id); err != nil {
		c.Log.Errorf("Failed to deleted tag %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success deleted tag")
	return nil
}
