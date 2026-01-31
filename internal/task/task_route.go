package task

import (
	"tasklybe/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterTaskRoutes(app *fiber.App) {
	task := app.Group("/task")
	task.Get("/", middleware.Auth(), HandleGetTasks)
	task.Get("/:id", middleware.Auth(), HandleGetDetailTask)
	task.Post("/", middleware.Auth(), HandleCreateTask)
	task.Put("/:id", middleware.Auth(), HandleEditTask)
	task.Delete("/:id", middleware.Auth(), HandleDeleteTask)
}
