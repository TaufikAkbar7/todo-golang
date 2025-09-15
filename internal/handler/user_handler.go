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

type UserHandler struct {
	Log     *logrus.Logger
	Service *service.UserService
}

func NewUserHandler(logger *logrus.Logger, service *service.UserService) *UserHandler {
	return &UserHandler{
		Log:     logger,
		Service: service,
	}
}

func (c *UserHandler) Register(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	request := new(model.UserCreateRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Errorf("Error validate request %v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}

	if err := c.Service.Create(newCtx, request); err != nil {
		if err == context.DeadlineExceeded {
			return ctx.Status(http.StatusGatewayTimeout).JSON(model.WebResponse[any]{Message: "operation timed out"})
		}
		if err.Error() == "Already registered user" {
			return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		c.Log.Errorf("Failed to register user : %+v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(model.WebResponse[any]{Message: err.Error()})
	}

	return ctx.JSON(model.WebResponse[any]{Message: "Success registered user"})
}

func (c *UserHandler) Login(ctx *fiber.Ctx) error {
	newCtx, cancel := context.WithTimeout(ctx.UserContext(), 3*time.Second)
	defer cancel()

	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Errorf("Failed to parse request body : %+v", err)
		return ctx.Status(http.StatusBadRequest).JSON(model.WebResponse[any]{Message: err.Error()})
	}

	response, err := c.Service.Login(newCtx, request)
	if err != nil {
		if err.Error() == "user not found" {
			return ctx.Status(http.StatusUnauthorized).JSON(model.WebResponse[any]{Message: err.Error()})
		}
		c.Log.Errorf("Failed to login user : %+v", err)
		return ctx.Status(http.StatusUnauthorized).JSON(model.WebResponse[any]{Message: "Wrong username or password"})
	}
	return ctx.JSON(model.WebResponse[*model.UserReponseLogin]{Data: response, Message: "Success login user"})
}
