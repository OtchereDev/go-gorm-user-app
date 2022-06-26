package controllers

import (
	"errors"
	"fmt"

	"github.com/OtchereDev/go-gorm-user-app/database"
	"github.com/OtchereDev/go-gorm-user-app/helpers"
	"github.com/OtchereDev/go-gorm-user-app/models"
	"github.com/OtchereDev/go-gorm-user-app/responses"
	"github.com/OtchereDev/go-gorm-user-app/serializers"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func UserSignUp(c *fiber.Ctx) error {
	userModel := database.ConnectToModel(models.UserModel{})

	var user serializers.UserSignUpSerializer

	// body parsing error check
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err}})
	}

	// body validation error checks
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadGateway, Message: "validation error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	// hash user password here
	hashedPassword, err := helpers.HashPassword(user.Password)

	if err != nil {
		return c.Status(400).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "Password Error", Data: &fiber.Map{"data": err.Error()}})
	}

	user.Password = string(hashedPassword)

	result := userModel.Create(&models.UserModel{Name: user.Name, Email: user.Email, Password: user.Password})

	// Checking for insertion error here

	if result.Error != nil {
		return c.Status(400).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": result.Error.Error()}})
	}

	var createdUser serializers.UserSignUpResponseSerializer
	// fetching the created user here
	if result.RowsAffected > 0 {
		userModel.First(&createdUser, "email = ?", user.Email)
	}

	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": createdUser}})
}

func UserLogin(c *fiber.Ctx) error {
	userModel := database.ConnectToModel(models.UserModel{})

	var body serializers.UserLoginSerializer

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if err := validate.Struct(&body); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {

			errorResponse := helpers.GenerateFromErrorMessage(ve)
			fmt.Println(errorResponse)
			return c.Status(fiber.StatusBadGateway).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Data: &fiber.Map{"data": errorResponse}})
		} else {
			return c.Status(fiber.StatusBadGateway).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Data: &fiber.Map{"data": err}})
		}
	}

	var user models.UserModel

	userModel.First(&user, "email = ?", body.Email)

	if (models.UserModel{}) == user {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Data: &fiber.Map{"data": &fiber.Map{"message": "Invalid credentials were provided"}}, Message: "error"})

	}

	if err := helpers.ComparePassword(user.Password, body.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.UserResponse{Status: fiber.StatusUnauthorized, Data: &fiber.Map{"data": &fiber.Map{"message": "Invalid credentials were provided"}}, Message: "error"})
	}

	t, err := helpers.GenerateJWT(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": &fiber.Map{"message": err.Error()}}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"token": t}})

}
