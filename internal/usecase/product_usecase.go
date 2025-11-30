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
	// if err := u.validateProduct(req); err != nil {
	// 	return nil, err
	// }

	product := &model.Product{
		Name:       strings.TrimSpace(req.Name),
		SKU:        strings.TrimSpace(req.SKU),
		CategoryID: int64(req.CategoryID),
		Description: func(desc *string) *string {
			if desc == nil {
				return nil
			}
			trimmed := strings.TrimSpace(*desc)
			return &trimmed
		}(req.Description),
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product %w", err)
	}

	return u.GetByID(ctx, product.ID)
}

// func (u *ProductUsecase) validateProduct(req *model.CreateProductRequest) error {
// 	if req.Name == "" || req.SKU == "" {
// 		return fmt.Errorf("invalid product")
// 	}
// 	return nil
// }

func (u *ProductUsecase) GetByID(ctx context.Context, id int64) (*model.Product, error) {

	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (u *ProductUsecase) GetAll(ctx context.Context) ([]*model.Product, error) {
	products, err := u.productRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	return products, nil
}

func (u *ProductUsecase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid product id")
	}

	if err := u.productRepo.DeleteByID(ctx, id); err != nil {
		return err
	}

	return nil
}
