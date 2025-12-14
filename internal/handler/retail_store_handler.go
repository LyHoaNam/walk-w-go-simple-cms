package handler

import (
	"simple-template/internal/usecase"
	"simple-template/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type RetailStoreHandler struct {
	RetailStoreUsecase *usecase.RetailStoreUsecase
}

func NewRetailStoreHandler(RetailStoreUsecase *usecase.RetailStoreUsecase) *RetailStoreHandler {
	return &RetailStoreHandler{
		RetailStoreUsecase: RetailStoreUsecase,
	}
}

func (h *RetailStoreHandler) GetAll(c *fiber.Ctx) error {
	RetailStores, err := h.RetailStoreUsecase.GetAll(c.Context())
	if err != nil {
		return response.BadRequest(c, "failed to get", err)
	}
	return response.Success(c, RetailStores, "Platform retrieved successfully")
}
