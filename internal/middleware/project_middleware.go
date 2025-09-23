package middleware

import (
	"golang-todo/internal/model"
	"golang-todo/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewProject(s *service.ProjectService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		project, err := s.GetByID(ctx.UserContext(), id)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(model.WebResponse[any]{Message: "Data not found"})
		}

		getUserId := GetUser(ctx).ID
		parseUUID, _ := uuid.Parse(getUserId)
		if project != nil && project.OwnerID == parseUUID {
			return ctx.Next()
		}
		return ctx.Status(http.StatusForbidden).JSON(model.WebResponse[any]{Message: "You are not authorized to manage this project"})
	}
}
