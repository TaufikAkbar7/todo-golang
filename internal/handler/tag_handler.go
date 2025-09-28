package handler

import (
	"context"
	"fmt"
	"golang-todo/internal/model"
	"golang-todo/internal/service"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TagHandler struct {
	Log       *logrus.Logger
	Service   *service.TagService
	Validator *validator.Validate
}

func NewTagHandler(logger *logrus.Logger, service *service.TagService, validator *validator.Validate) *TagHandler {
	return &TagHandler{
		Log:       logger,
		Service:   service,
		Validator: validator,
	}
}

func (c *TagHandler) GetAll(ctx *fiber.Ctx) error {
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

	// Map entity.Task to model.TagReponseGet
	var mappedResponse []model.TagReponseGet
	if response != nil {
		for _, tag := range *response {
			mappedResponse = append(mappedResponse, model.TagReponseGet{
				ID:        tag.ID,
				Name:      tag.Name,
				Color:     tag.Color,
				CreatedAt: tag.CreatedAt,
			})
		}
	}

	return ctx.JSON(model.WebResponse[*[]model.TagReponseGet]{Data: &mappedResponse, Message: "Success get tags"})
}

func (c *TagHandler) GetByID(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	id := ctx.Params("id")

	tag, err := c.Service.GetByID(newCtx, id)
	if err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		if err == fiber.ErrNotFound {
			return ctx.Status(http.StatusNotFound).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[*model.TagReponseGet]{Data: &model.TagReponseGet{
		ID:        tag.ID,
		Name:      tag.Name,
		Color:     tag.Color,
		CreatedAt: tag.CreatedAt,
	}, Message: "Success get tag"})
}

func (c *TagHandler) Create(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	request := new(model.TagCreateEditRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorf("Error validate request %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}

	if err := c.Service.Create(newCtx, request); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		fmt.Println(err)
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Message: "Success created tag"})
}

func (c *TagHandler) Update(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	id := ctx.Params("id")
	request := new(model.TagCreateEditRequest)
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
	return ctx.JSON(model.WebResponse[any]{Message: "Success updated tag"})
}

func (c *TagHandler) Delete(ctx *fiber.Ctx) error {
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
	return ctx.JSON(model.WebResponse[any]{Message: "Success deleted tag"})
}
