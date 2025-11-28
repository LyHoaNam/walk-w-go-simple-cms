package handler

import (
	"strconv"

	"simple-template/internal/model"
	"simple-template/internal/usecase"
	"simple-template/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// UserHandler xử lý các HTTP request liên quan đến user
type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

// NewUserHandler tạo instance mới của UserHandler
func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// CreateUser xử lý request tạo user mới
// POST /api/v1/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Parse request body
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	// Gọi usecase để tạo user
	user, err := h.userUsecase.CreateUser(c.Context(), &req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create user", err)
	}

	return response.Created(c, user, "User created successfully")
}

// GetUser xử lý request lấy thông tin user theo ID
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Parse ID từ URL params
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	// Gọi usecase để lấy user
	user, err := h.userUsecase.GetUserByID(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to get user", err)
	}

	return response.Success(c, user, "User retrieved successfully")
}

// GetAllUsers xử lý request lấy danh sách tất cả users
// GET /api/v1/users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Gọi usecase để lấy danh sách users
	users, err := h.userUsecase.GetAllUsers(c.Context())
	if err != nil {
		return response.InternalServerError(c, "Failed to get users", err)
	}

	return response.Success(c, users, "Users retrieved successfully")
}

// UpdateUser xử lý request cập nhật thông tin user
// PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Parse ID từ URL params
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

// DeleteUser xử lý request xóa user
// DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Parse ID từ URL params
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	// Gọi usecase để xóa user
	if err := h.userUsecase.DeleteUser(c.Context(), id); err != nil {
		if err.Error() == "user not found" {
			return response.NotFound(c, "User not found")
		}
		return response.InternalServerError(c, "Failed to delete user", err)
	}

	return response.Success(c, nil, "User deleted successfully")
}
