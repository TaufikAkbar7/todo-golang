package config

import (
	handler "golang-todo/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouteConfig struct {
	App                     *fiber.App
	AuthMiddleware          fiber.Handler
	ProjectMiddleware       fiber.Handler
	TaskMiddleware          fiber.Handler
	ProjectMemberMiddleware fiber.Handler
	UserHandler             *handler.UserHandler
	ProjectHandler          *handler.ProjectHandler
	TaskHandler             *handler.TaskHandler
	ProjectMemberHandler    *handler.ProjectMemberHandler
	TagHandler              *handler.TagHandler
}

func (c *RouteConfig) Setup() {
	// recover panic
	c.App.Use(recover.New())

	c.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	c.PublicRoutes()
	c.PrivateRoutes()
}

func (c *RouteConfig) PublicRoutes() {
	c.App.Post("/api/user/register", c.UserHandler.Register)
	c.App.Post("/api/user/login", c.UserHandler.Login)
}

func (c *RouteConfig) PrivateRoutes() {
	c.App.Use(c.AuthMiddleware)
	api := c.App.Group("/api")

	project := api.Group("/project")
	newProject := api.Group("/project/:projectId")
	project.Get("/", c.ProjectHandler.GetAll)
	project.Get(":id", c.ProjectHandler.GetByID)
	project.Post("/", c.ProjectHandler.Create)
	project.Patch(":id", c.ProjectMiddleware, c.ProjectHandler.Update)
	project.Delete(":id", c.ProjectMiddleware, c.ProjectHandler.Delete)
	project.Post(":id/invite", c.ProjectMemberMiddleware, c.ProjectMemberHandler.InviteUser)
	project.Delete(":id/delete/:userId", c.ProjectMemberMiddleware, c.ProjectMemberHandler.Delete)

	task := newProject.Group("/task")
	task.Get("/", c.TaskMiddleware, c.TaskHandler.GetAll)
	task.Get(":id", c.TaskMiddleware, c.TaskHandler.GetByID)
	task.Post("/", c.TaskMiddleware, c.TaskHandler.Create)
	task.Patch(":id", c.TaskMiddleware, c.TaskHandler.Update)
	task.Delete(":id", c.TaskMiddleware, c.TaskHandler.Delete)
	task.Post(":id/assign-tag", c.TaskHandler.AssignTag)
	task.Delete(":id/unassign-tag/:taskTagId", c.TaskHandler.UnassignTag)

	tag := newProject.Group("/tag")
	tag.Get("/", c.TaskMiddleware, c.TagHandler.GetAll)
	tag.Get(":id", c.TaskMiddleware, c.TagHandler.GetByID)
	tag.Post("/", c.TaskMiddleware, c.TagHandler.Create)
	tag.Patch(":id", c.TaskMiddleware, c.TagHandler.Update)
	tag.Delete(":id", c.TaskMiddleware, c.TagHandler.Delete)

}
