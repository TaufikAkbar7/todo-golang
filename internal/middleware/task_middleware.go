package middleware

import (
	"golang-todo/internal/model"
	"golang-todo/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func PermissionTask(ctx *fiber.Ctx, role int) bool {
	switch role {
	case 1, 2:
		return true
	case 3, 4:
		return ctx.Method() == fiber.MethodGet
	default:
		return false
	}
}

func NewTask(s *service.ProjectMemberService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		roleMember := GetRoleMember(ctx)
		if roleMember != nil {
			if isGranted := PermissionTask(ctx, roleMember.RoleID); isGranted {
				return ctx.Next()
			}
			return ctx.Status(http.StatusForbidden).JSON(model.WebResponse[any]{Message: "You are not authorized to manage task on this project"})
		}

		userID := GetUser(ctx).ID
		projectID := ctx.Params("projectId")
		newProjectID, _ := uuid.Parse(projectID)

		projectMember, err := s.GetByID(ctx.UserContext(), userID, newProjectID)
		if err != nil {
			return ctx.Status(http.StatusForbidden).JSON(model.WebResponse[any]{Message: "You are not authorized to manage task on this project"})
		}

		isGranted := PermissionTask(ctx, projectMember.RoleID)
		if isGranted {
			ctx.Locals("role_member", &model.ProjectMemberLocal{
				UserID:    userID,
				RoleID:    projectMember.RoleID,
				ProjectID: newProjectID,
			})
			return ctx.Next()
		}
		return ctx.Status(http.StatusForbidden).JSON(model.WebResponse[any]{Message: "You are not authorized to manage task on this project"})
	}
}
