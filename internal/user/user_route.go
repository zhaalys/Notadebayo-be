package user

import "github.com/gofiber/fiber/v2"

func RegisterUserRoutes(app *fiber.App) {
	user := app.Group("/user")
	user.Post("/register", HandleRegister)
	user.Post("/login", HandleLogin)
}
