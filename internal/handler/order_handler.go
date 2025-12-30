package handler

import (
	"simple-template/internal/model"
	"simple-template/internal/usecase"
	"simple-template/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderUsecase *usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var req model.CreateOrders
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", err)
	}
	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, "validation failed", err)
	}
	orders, err := h.orderUsecase.CreateOrders(c.Context(), &req)
	if err != nil {
		return response.BadRequest(c, "create failed", err)
	}

	return response.Success(c, orders, "success")
}

func (h *OrderHandler) GetAll(c *fiber.Ctx) error {
	orders, err := h.orderUsecase.GetOrdersPage(c.Context())
	if err != nil {
		return response.BadRequest(c, "Failed to fetch", err)
	}
	return response.Success(c, orders, "orders retrieved successfully")
}
