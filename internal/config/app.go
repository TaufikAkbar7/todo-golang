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

	tagRepo := repository.NewTagRepository(config.DB, config.Log)
	tagService := service.NewTagService(config.DB, tagRepo, config.Log)
	tagHandler := handler.NewTagHandler(config.Log, tagService, config.Validate)

	taskRepo := repository.NewTaskRepository(config.DB, config.Log)
	taskService := service.NewTaskService(config.DB, taskRepo, config.Log, tagRepo)
	taskHandler := handler.NewTaskHandler(config.Log, taskService, config.Validate)

	authMiddleware := middleware.NewAuth(userService)
	projectMiddleware := middleware.NewProject(projectService)
	taskMiddleware := middleware.NewTask(projectMemberService)
	projectMemberMiddleware := middleware.NewProjectMember(projectMemberService)

	routeConfig := RouteConfig{
		App:                     config.App,
		AuthMiddleware:          authMiddleware,
		ProjectMiddleware:       projectMiddleware,
		TaskMiddleware:          taskMiddleware,
		ProjectMemberMiddleware: projectMemberMiddleware,
		UserHandler:             userHandler,
		ProjectHandler:          projectHandler,
		TaskHandler:             taskHandler,
		ProjectMemberHandler:    projectMemberHandler,
		TagHandler:              tagHandler,
	}
	routeConfig.Setup()
}
