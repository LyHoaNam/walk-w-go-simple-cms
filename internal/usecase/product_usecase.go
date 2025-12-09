package usecase

import (
	"context"
	"fmt"
	"simple-template/internal/model"
	"simple-template/internal/repository"
	"simple-template/internal/utils"
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

	if len(req.Variants) <= 0 {
		return nil, fmt.Errorf("invalid request")
	}

	product := &model.Product{
		Name:        strings.TrimSpace(req.Name),
		SKU:         strings.TrimSpace(req.SKU),
		Status:      req.Status,
		ImgUrl:      strings.TrimSpace(req.ImgUrl),
		CategoryID:  req.CategoryID,
		Description: utils.TrimStringPointer(req.Description),
		Dimension:   utils.TrimStringPointer(req.Dimension),
		Weight:      req.Weight,
		Brand:       utils.TrimStringPointer(req.Brand),
		Material:    utils.TrimStringPointer(req.Material),
		Origin:      utils.TrimStringPointer(req.Origin),
	}

	for _, reqVariant := range req.Variants {
		variant := &model.ProductVariant{
			Name:        strings.TrimSpace(reqVariant.Name),
			DisplayName: strings.TrimSpace(reqVariant.DisplayName),
			IsRequire: func(requitable *int16) int16 {
				if requitable == nil {
					return int16(0)
				}
				return *requitable
			}(reqVariant.IsRequire),
			DisplayOrder: func(order *int64) int64 {
				if order == nil {
					return int64(0)
				}
				return *order
			}(reqVariant.DisplayOrder),
		}

		if len(reqVariant.Values) > 0 {
			for _, reqValue := range reqVariant.Values {
				value := &model.ProductVariantValue{
					Value:         strings.TrimSpace(reqValue.Value),
					DisplayOrder:  reqValue.DisplayOrder,
					StockQuantity: reqValue.StockQuantity,
				}

				variant.Values = append(variant.Values, *value)
			}
		}
		product.Variant = append(product.Variant, *variant)
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
