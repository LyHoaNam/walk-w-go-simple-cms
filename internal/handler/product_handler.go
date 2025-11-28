package handler

import (
	"simple-template/internal/model"
	"simple-template/internal/usecase"
	"simple-template/pkg/response"

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

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {

	var req model.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	product, err := h.productUsecase.CreateProduct(c.Context(), &req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create product", err)
	}

	return response.Created(c, product, "Product created successfully")
}
