package middleware

import (
	"golang-todo/internal/model"
	"golang-todo/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewProjectMember(s *service.ProjectMemberService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if GetRoleMember(ctx) != nil {
			projectMember := GetRoleMember(ctx)
			if projectMember.RoleID == 1 || projectMember.RoleID == 2 {
				return ctx.Next()
			}
			return ctx.Status(http.StatusForbidden).JSON(model.WebResponse[any]{Message: "You are not authorized to manage this project"})
		}

		projectID := ctx.Params("id")
		userID := GetUser(ctx).ID
		newProjectID, _ := uuid.Parse(projectID)
		projectMember, err := s.GetByID(ctx.UserContext(), userID, newProjectID)
		if err != nil {
			return ctx.Status(http.StatusForbidden).JSON(model.WebResponse[any]{Message: "You are not authorized to manage this project"})
		}

		if projectMember.RoleID == 1 || projectMember.RoleID == 2 {
			ctx.Locals("role_member", &model.ProjectMemberLocal{
				UserID:    userID,
				RoleID:    projectMember.RoleID,
				ProjectID: newProjectID,
			})
			return ctx.Next()
		}
		return ctx.Status(http.StatusForbidden).JSON(model.WebResponse[any]{Message: "You are not authorized to manage this project"})
	}
}

func GetRoleMember(ctx *fiber.Ctx) *model.ProjectMemberLocal {
	role, ok := ctx.Locals("role_member").(*model.ProjectMemberLocal)
	if !ok {
		return nil
	}
	return role
}
