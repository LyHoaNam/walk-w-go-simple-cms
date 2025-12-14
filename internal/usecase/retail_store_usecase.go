package usecase

import (
	"context"
	"simple-template/internal/model"
	"simple-template/internal/repository"
)

type RetailStoreUsecase struct {
	RetailStoreRepo *repository.RetailStoreRepository
}

func NewRetailStoreUsecase(RetailStoreR *repository.RetailStoreRepository) *RetailStoreUsecase {
	return &RetailStoreUsecase{
		RetailStoreRepo: RetailStoreR,
	}
}

func (r *RetailStoreUsecase) GetAll(ctx context.Context) ([]*model.RetailStore, error) {
	RetailStores, err := r.RetailStoreRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return RetailStores, nil
}
