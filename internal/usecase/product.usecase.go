package usecase

import (
	"context"
	"fmt"
	"simple-template/internal/model"
	"simple-template/internal/repository"
	"strings"
)

type ProductUsecase struct {
	productRepo *repository.ProductRepository
}

func NewProductUsecase(productRepo *repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {

	if err := u.ValidateProduct(req); err != nil {
		return nil, err
	}
	product := &model.Product{
		Name:       strings.TrimSpace(req.Name),
		Price:      req.Price,
		CategoryID: req.CategoryID,
	}
	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	createdProduct, err := u.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created product: %w", err)
	}
	return createdProduct, nil
}

func (u *ProductUsecase) ValidateProduct(req *model.CreateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("invalid product name")
	}
	if req.Price <= 0 {
		return fmt.Errorf("invalid product price")
	}
	return nil
}
