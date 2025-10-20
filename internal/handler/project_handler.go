package handler

import (
	"context"
	"golang-todo/internal/middleware"
	"golang-todo/internal/model"
	"golang-todo/internal/service"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProjectHandler struct {
	Log       *logrus.Logger
	Service   *service.ProjectService
	Validator *validator.Validate
}

func NewProjectHandler(logger *logrus.Logger, service *service.ProjectService, validator *validator.Validate) *ProjectHandler {
	return &ProjectHandler{
		Log:       logger,
		Service:   service,
		Validator: validator,
	}
}

func (c *ProjectHandler) GetAll(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	userID := middleware.GetUser(ctx).ID
	response, err := c.Service.GetAll(newCtx, userID)
	if err != nil {
		emptyData := []any{}
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[[]any]{Data: emptyData, Message: "operation timed out"})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[[]any]{Data: emptyData, Message: "Internal server error"})
	}
	return ctx.JSON(model.WebResponse[*[]model.ProjectReponseGet]{Data: response, Message: "Success get all projects"})
}

func (c *ProjectHandler) GetByID(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	id := ctx.Params("id")
	projectID, _ := uuid.Parse(id)
	userID := middleware.GetUser(ctx).ID

	project, err := c.Service.GetByID(newCtx, projectID, userID)
	if err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		if err == fiber.ErrNotFound {
			return ctx.Status(http.StatusNotFound).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[*model.ProjectReponseGet]{Data: &model.ProjectReponseGet{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID,
		CreatedAt:   project.CreatedAt,
	}, Message: "Success get project"})
}

func (c *ProjectHandler) Create(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	request := new(model.ProjectCreateEditRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorf("Error validate request %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Errorf("Invalid request body  : %+v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	// get id user login
	id := middleware.GetUser(ctx).ID

	if err := c.Service.Create(newCtx, request, id); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Message: "Success create project"})
}

func (c *ProjectHandler) Update(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	id := ctx.Params("id")
	projectID, _ := uuid.Parse(id)
	userID := middleware.GetUser(ctx).ID
	request := new(model.ProjectCreateEditRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorf("Error validate request %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	if err := c.Validator.Struct(request); err != nil {
		c.Log.Errorf("Invalid request body  : %+v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}

	if err := c.Service.Update(newCtx, request, projectID, userID); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		if err == fiber.ErrNotFound {
			return ctx.Status(http.StatusNotFound).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Message: "Success update project"})
}

func (c *ProjectHandler) Delete(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	id := ctx.Params("id")
	projectID, _ := uuid.Parse(id)
	userID := middleware.GetUser(ctx).ID

	if err := c.Service.Delete(newCtx, projectID, userID); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		if err == fiber.ErrNotFound {
			return ctx.Status(http.StatusNotFound).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Message: "Success deleted project"})
}
