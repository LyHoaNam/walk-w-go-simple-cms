package handler

import (
	"simple-template/internal/usecase"
	"simple-template/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type PaymentMethodsHandler struct {
	PaymentMethodsUsecase *usecase.PaymentMethodsUsecase
}

func NewPaymentMethodsHandler(PaymentMethodsUsecase *usecase.PaymentMethodsUsecase) *PaymentMethodsHandler {
	return &PaymentMethodsHandler{
		PaymentMethodsUsecase: PaymentMethodsUsecase,
	}
}

func (h *PaymentMethodsHandler) GetAll(c *fiber.Ctx) error {
	PaymentMethods, err := h.PaymentMethodsUsecase.GetAll(c.Context())
	if err != nil {
		return response.BadRequest(c, "failed to get", err)
	}
	return response.Success(c, PaymentMethods, "payment methods retrieved successfully")
}
