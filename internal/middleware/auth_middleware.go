package middleware

import (
	"golang-todo/internal/model"
	"golang-todo/internal/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userService *service.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userService.Log.Debugf("Authorization : %s", request.Token)

		// check if format "Bearer <token>"
		parts := strings.Split(request.Token, " ")
		if len(parts) > 2 && parts[0] != "Bearer" {
			return ctx.Status(http.StatusUnauthorized).JSON(model.WebResponse[any]{Message: "Missing or malformed JWT"})
		} else if len(parts) == 2 && parts[0] == "Bearer" {
			request.Token = parts[1]
		}

		auth, err := userService.Verify(ctx.UserContext(), request)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				return ctx.Status(http.StatusUnauthorized).JSON(model.WebResponse[any]{Message: "Token is expired"})
			}
			return ctx.Status(http.StatusUnauthorized).JSON(model.WebResponse[any]{Message: "Unauthorized"})
		}

		userService.Log.Debugf("User login: %+v", auth.Username)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
