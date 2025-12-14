package usecase

import (
	"context"
	"simple-template/internal/model"
	"simple-template/internal/repository"
)

type PaymentMethodsUsecase struct {
	PaymentMethodsRepo *repository.PaymentMethodsRepository
}

func NewPaymentMethodsUsecase(PaymentMethodsR *repository.PaymentMethodsRepository) *PaymentMethodsUsecase {
	return &PaymentMethodsUsecase{
		PaymentMethodsRepo: PaymentMethodsR,
	}
}

func (r *PaymentMethodsUsecase) GetAll(ctx context.Context) ([]*model.PaymentMethods, error) {
	PaymentMethods, err := r.PaymentMethodsRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return PaymentMethods, nil
}
