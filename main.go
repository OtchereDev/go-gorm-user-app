package main

import (
	"log"

	"github.com/OtchereDev/go-gorm-user-app/database"
	"github.com/OtchereDev/go-gorm-user-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	database.ConnectDB()

	app := fiber.New()

	routes.UserRoutes(app)
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to My Api")
	})

	app.Listen(":3000")
}
