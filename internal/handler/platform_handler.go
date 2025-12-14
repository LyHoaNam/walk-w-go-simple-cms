package handler

import (
	"simple-template/internal/usecase"
	"simple-template/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type PlatformHandler struct {
	platformUsecase *usecase.PlatformUsecase
}

func NewPlatformHanlder(platformUsecase *usecase.PlatformUsecase) *PlatformHandler {
	return &PlatformHandler{
		platformUsecase: platformUsecase,
	}
}

func (h *PlatformHandler) GetAll(c *fiber.Ctx) error {
	platforms, err := h.platformUsecase.GetAll(c.Context())
	if err != nil {
		return response.BadRequest(c, "failed to get", err)
	}
	return response.Success(c, platforms, "Plat form retrieved successfully")
}
