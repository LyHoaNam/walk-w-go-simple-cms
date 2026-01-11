package handler

import (
	"simple-template/internal/model"
	"simple-template/internal/usecase"
	"simple-template/pkg/pagination"
	"simple-template/pkg/response"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
}

var validate = validator.New()

func NewProductHandler(productUsecase *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

// GET /api/v1/product/:id
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid product ID", err)
	}
	product, err := h.productUsecase.GetByID(c.Context(), id)
	if err != nil {
		return response.BadRequest(c, "Failed to fetch", err)
	}
	return response.Success(c, product, "product retrieved successfully")
}

// GET /api/v1/products
func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	var req pagination.Request
	if err := c.QueryParser(&req); err != nil {
		return response.BadRequest(c, "invalid query parameters", err)
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, "validation failed", err)
	}

	products, err := h.productUsecase.GetAll(c.Context(), &req)
	if err != nil {
		return response.BadRequest(c, "Failed to fetch", err)
	}
	return response.Success(c, products, "product retrieved successfully")
}

// POST /api/v1/product
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req model.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body", err)
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, "validation failed", err)
	}

	product, err := h.productUsecase.CreateProduct(c.Context(), &req)
	if err != nil {
		return h.handleCreateError(c, err)
	}
	return response.Created(c, product, "product created successfully")
}

// DELETE /api/v1/product
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	if err := h.productUsecase.Delete(c.Context(), id); err != nil {
		if err.Error() == "product not found" {
			return response.InternalServerError(c, "Failed to delete user", err)
		}
	}

	return response.Success(c, nil, "Product deleted successfully")
}

func (h *ProductHandler) handleCreateError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// Business validation errors
	if strings.Contains(errMsg, "already exists") ||
		strings.Contains(errMsg, "duplicate") {
		return response.BadRequest(c, errMsg, err)
	}

	if strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "required") ||
		strings.Contains(errMsg, "cannot be empty") ||
		strings.Contains(errMsg, "must have") {
		return response.BadRequest(c, errMsg, err)
	}

	if strings.Contains(errMsg, "not found") {
		return response.NotFound(c, errMsg)
	}

	// Unexpected errors
	return response.InternalServerError(c, "failed to create product", err)
}
