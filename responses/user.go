package responses

import "github.com/gofiber/fiber/v2"

type UserResponse struct {
	Status  int16      `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}
