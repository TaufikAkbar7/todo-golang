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
	userService := service.NewUserService(config.DB, userRepo, config.Log, config.Validate)
	userHandler := handler.NewUserHandler(config.Log, userService)

	projectRepo := repository.NewProjectRepository(config.DB, config.Log)
	projectService := service.NewProjectService(config.DB, projectRepo, config.Log, config.Validate)
	projectHandler := handler.NewProjectHandler(config.Log, projectService)

	taskRepo := repository.NewTaskRepository(config.DB, config.Log)
	taskService := service.NewTaskService(config.DB, taskRepo, config.Log, config.Validate)
	taskHandler := handler.NewTaskHandler(config.Log, taskService)

	authMiddleware := middleware.NewAuth(userService)
	routeConfig := RouteConfig{
		App:            config.App,
		AuthMiddleware: authMiddleware,
		UserHandler:    userHandler,
		ProjectHandler: projectHandler,
		TaskHandler:    taskHandler,
	}
	routeConfig.Setup()
}
