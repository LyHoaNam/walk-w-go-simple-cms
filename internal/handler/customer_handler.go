package handler

import (
	"simple-template/internal/model"
	"simple-template/internal/usecase"
	"simple-template/pkg/response"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	customerUsecase *usecase.CustomerUsecase
}

func NewCustomerHandler(customerUsecase *usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: customerUsecase,
	}
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var req model.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}
	if err := validator.New().Struct(req); err != nil {
		return response.BadRequest(c, "validation failed", err)
	}

	customer, err := h.customerUsecase.CreateCustomer(c.Context(), &req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create user", err)
	}
	return response.Created(c, customer, "Customer created successfully")

}

func (h *CustomerHandler) GetCustomer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id", err)
	}

	customer, err := h.customerUsecase.GetCustomerByID(c.Context(), id)
	if err != nil {
		return response.BadRequest(c, "customer not found", err)
	}
	return response.Success(c, customer, "Customer retrieved successfully")
}

func (h *CustomerHandler) GetAllCustomers(c *fiber.Ctx) error {

	customers, err := h.customerUsecase.GetAllCustomer(c.Context())
	if err != nil {
		return response.BadRequest(c, "customer not found", err)
	}
	return response.Success(c, customers, "Customer retrieved successfully")
}

func (h *CustomerHandler) UpdateCustomers(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id", err)
	}

	var req model.UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body", err)
	}

	customers, err := h.customerUsecase.UpdateCustomer(c.Context(), id, &req)
	if err != nil {
		return response.BadRequest(c, "customer not found", err)
	}
	return response.Success(c, customers, "Customer updated successfully")
}

func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.BadRequest(c, "invalid id", err)
	}

	if err := h.customerUsecase.DeleteCustomer(c.Context(), id); err != nil {
		return response.BadRequest(c, "failed", err)
	}
	return nil
}
