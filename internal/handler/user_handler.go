package handler

import (
	"strconv"

	"simple-template/internal/model"
	"simple-template/internal/usecase"
	"simple-template/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles all HTTP requests related to users
type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// CreateUser handles the request to create a new user
// POST /api/v1/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Parse request body
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	// Call usecase to create user
	user, err := h.userUsecase.CreateUser(c.Context(), &req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create user", err)
	}

	return response.Created(c, user, "User created successfully")
}

// GetUser handles the request to get user information by ID
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Parse ID from URL params
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	// Call usecase to get user
	user, err := h.userUsecase.GetUserByID(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to get user", err)
	}

	return response.Success(c, user, "User retrieved successfully")
}

// GetAllUsers handles the request to get a list of all users
// GET /api/v1/users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Call usecase to get list of users
	users, err := h.userUsecase.GetAllUsers(c.Context())
	if err != nil {
		return response.InternalServerError(c, "Failed to get users", err)
	}

	return response.Success(c, users, "Users retrieved successfully")
}

// UpdateUser handles the request to update user information
// PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Parse ID from URL params
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	// Parse request body
	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	// Gọi usecase để update user
	user, err := h.userUsecase.UpdateUser(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to update user", err)
	}

	return response.Success(c, user, "User updated successfully")
}

// DeleteUser handles the request to delete a user
// DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Parse ID from URL params
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	// Call usecase to delete user
	if err := h.userUsecase.DeleteUser(c.Context(), id); err != nil {
		if err.Error() == "user not found" {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to delete user", err)
	}

	return response.Success(c, nil, "User deleted successfully")
}
