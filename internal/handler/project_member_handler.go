package handler

import (
	"context"
	"golang-todo/internal/model"
	"golang-todo/internal/service"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProjectMemberHandler struct {
	Log     *logrus.Logger
	Service *service.ProjectMemberService
}

func NewProjectMemberHandler(logger *logrus.Logger, service *service.ProjectMemberService) *ProjectMemberHandler {
	return &ProjectMemberHandler{
		Log:     logger,
		Service: service,
	}
}

func (c *ProjectMemberHandler) InviteUser(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	request := new(model.ProjectMemberCreateEditRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorf("Error validate request %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}

	id := ctx.Params("id")
	parseUUID, _ := uuid.Parse(id)
	request.ProjectID = parseUUID
	if err := c.Service.Create(newCtx, request); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		c.Log.Debug(err)
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Message: "Success invited user"})
}

func (c *ProjectMemberHandler) Delete(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	id := ctx.Params("id")
	userId := ctx.Params("userId")

	if err := c.Service.Delete(newCtx, id, userId); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		if err == fiber.ErrNotFound {
			return ctx.Status(http.StatusNotFound).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}
	return ctx.JSON(model.WebResponse[any]{Message: "Success deleted member"})
}
