package service

import (
	"context"
	"database/sql"
	"errors"
	"golang-todo/internal/entity"
	"golang-todo/internal/helper"
	"golang-todo/internal/model"
	"golang-todo/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TaskService struct {
	DB   *sqlx.DB
	Repo *repository.TaskRepository
	Log  *logrus.Logger
}

func NewTaskService(db *sqlx.DB, repo *repository.TaskRepository, log *logrus.Logger) *TaskService {
	return &TaskService{
		DB:   db,
		Repo: repo,
		Log:  log,
	}
}

func (c *TaskService) GetAll(ctx context.Context) (*[]entity.Task, error) {
	tasks, err := c.Repo.GetAll(ctx)
	if err != nil {
		c.Log.Errorf("Failed to get all data tasks %v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Info("Success get all tasks")
	return tasks, nil
}

func (c *TaskService) GetByID(ctx context.Context, id string) (*entity.Task, error) {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	task, err := c.Repo.GetByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			c.Log.Info("Data not found")
			return nil, fiber.ErrNotFound
		}
		c.Log.Errorf("Failed to get data task %v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Info("Success get task by id")
	return task, nil
}

func (c *TaskService) Create(ctx context.Context, req *model.TaskCreateEditRequest) error {
	newID, _ := uuid.NewV7()
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	var dueDatePointer *time.Time
	// check if the due date was provided in the request
	if req.DueDate != nil && *req.DueDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid due_date format, please use RFC3339")
		}
		dueDatePointer = &parsedTime
	}

	dateNow := helper.GetDateNow()
	payload := &entity.Task{
		ID:        newID,
		Title:     req.Title,
		ProjectID: req.ProjectID,
		CreatedAt: dateNow,
		UpdatedAt: dateNow,
		DueDate:   dueDatePointer,
		Status:    req.Status,
		Priority:  req.Priority,
	}

	if err := c.Repo.Create(ctx, tx, payload); err != nil {
		c.Log.Errorf("Failed create task %v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success created task")
	return nil
}

func (c *TaskService) Update(ctx context.Context, req *model.TaskCreateEditRequest, id string) error {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	task, err := c.Repo.GetByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			c.Log.Info("Data not found")
			return fiber.ErrNotFound
		}
		c.Log.Errorf("Failed to get data task %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Debug(req.DueDate, "due_Date")
	var dueDatePointer *time.Time

	// check if the due date was provided in the request
	if req.DueDate != nil && *req.DueDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid due_date format, please use RFC3339")
		}
		dueDatePointer = &parsedTime
	}

	dateNow := helper.GetDateNow()
	payload := &entity.Task{
		ID:        task.ID,
		Title:     req.Title,
		ProjectID: req.ProjectID,
		UpdatedAt: dateNow,
		DueDate:   dueDatePointer,
		Status:    req.Status,
		Priority:  req.Priority,
	}

	if err := c.Repo.Update(ctx, tx, payload); err != nil {
		c.Log.Errorf("Failed create task %v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success updated task")
	return nil
}

func (c *TaskService) Delete(ctx context.Context, id string) error {
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
		c.Log.Errorf("Failed to get data task %v", err)
		return fiber.ErrInternalServerError
	}

	if err := c.Repo.Delete(ctx, tx, id); err != nil {
		c.Log.Errorf("Failed to deleted user %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Info("Success deleted task")
	return nil
}
