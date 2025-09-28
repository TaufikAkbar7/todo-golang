package config

import (
	handler "golang-todo/internal/handler"
	"golang-todo/internal/middleware"
	"golang-todo/internal/repository"
	"golang-todo/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	DB       *sqlx.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	userRepo := repository.NewUserRepository(config.DB, config.Log)
	userService := service.NewUserService(config.DB, userRepo, config.Log)
	userHandler := handler.NewUserHandler(config.Log, userService, config.Validate)

	projectMemberRepo := repository.NewProjectMemberRepository(config.DB, config.Log)
	projectMemberService := service.NewProjectMemberService(config.DB, projectMemberRepo, config.Log)
	projectMemberHandler := handler.NewProjectMemberHandler(config.Log, projectMemberService, config.Validate)

	projectRepo := repository.NewProjectRepository(config.DB, config.Log)
	projectService := service.NewProjectService(config.DB, projectRepo, config.Log, projectMemberRepo)
	projectHandler := handler.NewProjectHandler(config.Log, projectService, config.Validate)

	taskRepo := repository.NewTaskRepository(config.DB, config.Log)
	taskService := service.NewTaskService(config.DB, taskRepo, config.Log)
	taskHandler := handler.NewTaskHandler(config.Log, taskService, config.Validate)

	authMiddleware := middleware.NewAuth(userService)
	projectMiddleware := middleware.NewProject(projectService)
	routeConfig := RouteConfig{
		App:                  config.App,
		AuthMiddleware:       authMiddleware,
		ProjectMiddleware:    projectMiddleware,
		UserHandler:          userHandler,
		ProjectHandler:       projectHandler,
		TaskHandler:          taskHandler,
		ProjectMemberHandler: projectMemberHandler,
	}
	routeConfig.Setup()
}
