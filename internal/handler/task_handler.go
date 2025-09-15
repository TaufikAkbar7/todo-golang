package handler

import (
	"context"
	"golang-todo/internal/model"
	"golang-todo/internal/service"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TaskHandler struct {
	Log     *logrus.Logger
	Service *service.TaskService
}

func NewTaskHandler(logger *logrus.Logger, service *service.TaskService) *TaskHandler {
	return &TaskHandler{
		Log:     logger,
		Service: service,
	}
}

func (c *TaskHandler) GetAll(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	response, err := c.Service.GetAll(newCtx)
	if err != nil {
		emptyData := []any{}
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[[]any]{Data: emptyData, Message: "operation timed out"})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[[]any]{Data: emptyData, Message: "Internal server error"})
	}

	// Map entity.Task to model.TaskReponseGet
	var mappedResponse []model.TaskReponseGet
	if response != nil {
		for _, task := range *response {
			mappedResponse = append(mappedResponse, model.TaskReponseGet{
				ID:        task.ID,
				Title:     task.Title,
				ProjectID: task.ProjectID,
				CreatedAt: task.CreatedAt,
			})
		}
	}

	return ctx.JSON(model.WebResponse[*[]model.TaskReponseGet]{Data: &mappedResponse, Message: "Success get tasks"})
}

func (c *TaskHandler) GetByID(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	id := ctx.Params("id")

	task, err := c.Service.GetByID(newCtx, id)
	if err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		if err == fiber.ErrNotFound {
			return ctx.Status(http.StatusNotFound).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[*model.TaskReponseGet]{Data: &model.TaskReponseGet{
		ID:        task.ID,
		Title:     task.Title,
		ProjectID: task.ProjectID,
		CreatedAt: task.CreatedAt,
	}, Message: "Success get task"})
}

func (c *TaskHandler) Create(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	request := new(model.TaskCreateEditRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorf("Error validate request %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}

	if err := c.Service.Create(newCtx, request); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Message: "Success created task"})
}

func (c *TaskHandler) Update(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	id := ctx.Params("id")
	request := new(model.TaskCreateEditRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorf("Error validate request %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}

	if err := c.Service.Update(newCtx, request, id); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		if err == fiber.ErrNotFound {
			return ctx.Status(http.StatusNotFound).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Message: "Success updated task"})
}

func (c *TaskHandler) Delete(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	id := ctx.Params("id")

	if err := c.Service.Delete(newCtx, id); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		if err == fiber.ErrNotFound {
			return ctx.Status(http.StatusNotFound).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Message: "Success deleted task"})
}
