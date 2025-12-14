package usecase

import (
	"context"
	"simple-template/internal/model"
	"simple-template/internal/repository"
)

type PlatformUsecase struct {
	platformRepo *repository.PlatformRepository
}

func NewPlatformUsecase(platformR *repository.PlatformRepository) *PlatformUsecase {
	return &PlatformUsecase{
		platformRepo: platformR,
	}
}

func (r *PlatformUsecase) GetAll(ctx context.Context) ([]*model.Platform, error) {
	platforms, err := r.platformRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return platforms, nil
}
