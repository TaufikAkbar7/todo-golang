package config

import (
	handler "golang-todo/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
	UserHandler    *handler.UserHandler
	ProjectHandler *handler.ProjectHandler
	TaskHandler    *handler.TaskHandler
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
	c.App.Post("/api/users/register", c.UserHandler.Register)
	c.App.Post("/api/users/login", c.UserHandler.Login)
}

func (c *RouteConfig) PrivateRoutes() {
	c.App.Use(c.AuthMiddleware)
	api := c.App.Group("/api")

	projects := api.Group("/projects")
	projects.Get("/", c.ProjectHandler.GetAll)
	projects.Get(":id", c.ProjectHandler.GetByID)
	projects.Post("/", c.ProjectHandler.Create)
	projects.Patch(":id", c.ProjectHandler.Update)
	projects.Delete(":id", c.ProjectHandler.Delete)

	tasks := api.Group("/tasks")
	tasks.Get("/", c.TaskHandler.GetAll)
	tasks.Get(":id", c.TaskHandler.GetByID)
	tasks.Post("/", c.TaskHandler.Create)
	tasks.Patch(":id", c.TaskHandler.Update)
	tasks.Delete(":id", c.TaskHandler.Delete)
}
