package user

import (
	"tasklybe/internal/dto"
	"tasklybe/internal/validation"

	"github.com/gofiber/fiber/v2"
)

func HandleRegister(c *fiber.Ctx) error {
	var req RegisterUserRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseWrapper[User]{
			Data:    nil,
			Success: false,
			Message: "Failed! Validation error.",
			Error:   validation.FormatValidationError(err),
		})
	}

	u, err := RegisterUser(req)
	if err != nil {
		if err == ErrEmailAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(dto.ResponseWrapper[User]{
				Data:    nil,
				Success: false,
				Message: "Failed! Email already exists.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[User]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[UserResponse]{
		Data:    ToUserResponse(u),
		Success: true,
		Message: "Success! User registered.",
	})
}

func HandleLogin(c *fiber.Ctx) error {
	var req LoginUserRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseWrapper[User]{
			Data:    nil,
			Success: false,
			Message: "Failed! Validation error.",
			Error:   validation.FormatValidationError(err),
		})
	}

	u, apiKey, err := LoginUser(req)
	if err != nil {
		if err == ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ResponseWrapper[User]{
				Data:    nil,
				Success: false,
				Message: "Failed! Invalid credentials.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ResponseWrapper[User]{
			Data:    nil,
			Success: false,
			Message: "Failed! Something went wrong.",
		})
	}

	return c.JSON(dto.ResponseWrapper[UserResponse]{
		Data: &UserResponse{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			Email:     u.Email,
			Name:      u.Name,
			ApiKey:    apiKey,
		},
		Success: true,
		Message: "Success! Login success.",
	})
}
