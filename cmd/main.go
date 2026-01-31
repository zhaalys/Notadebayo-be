package main

import (
	"log"
	"tasklybe/internal/task"
	"tasklybe/internal/user"
	"tasklybe/pkg/db"

	_ "tasklybe/docs" // import docs

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// @title Taskly API
// @version 1.0
// @description Taskly backend API
// @host localhost:3002
// @BasePath /
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	app := fiber.New() //initialize web server

	// Setup CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	log.Println("Server started successfully!")
	db.Connect()

	log.Println("Migrating table")
	if err := db.DB.AutoMigrate(&task.Task{}); err != nil {
		log.Fatal("Failed to migrate table", err)
	}
	if err := db.DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("Failed to migrate table", err)
	}

	task.RegisterTaskRoutes(app)
	user.RegisterUserRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Static("/openapi.json", "./docs/swagger.json")
	app.Static("/docs", "./public")

	log.Fatal(app.Listen(":3002"))
}
