package routes

import (
	"github.com/OtchereDev/go-gorm-user-app/controllers"
	"github.com/OtchereDev/go-gorm-user-app/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	app.Post("/user", controllers.UserSignUp)
	app.Post("/login", controllers.UserLogin)

	app.Use(middlewares.New(middlewares.Config{}))

	app.Get("/protected", func(c *fiber.Ctx) error {
		claimData := c.Locals("user")

		if claimData == nil {
			return c.SendString("Jwt was bypassed")
		}
		return c.JSON(claimData)

	})
}
