package helpers

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return "There was an error with this field" // default error
}

func GenerateFromErrorMessage(ve validator.ValidationErrors) []fiber.Map {

	out := make([]fiber.Map, len(ve))
	for i, fe := range ve {
		out[i] = fiber.Map{fe.Field(): msgForTag(fe)}
	}
	return out

}
