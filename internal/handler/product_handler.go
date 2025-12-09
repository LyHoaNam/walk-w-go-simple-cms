package handler

import (
	"simple-template/internal/model"
	"simple-template/internal/usecase"
	"simple-template/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
}

func NewProductHandler(productUsecase *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

// POST /api/v1/product/:id
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
	products, err := h.productUsecase.GetAll(c.Context())
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
	product, err := h.productUsecase.CreateProduct(c.Context(), &req)
	if err != nil {
		return response.InternalServerError(c, "failed to create product", err)
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
