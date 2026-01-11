package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"simple-template/internal/model"
	"simple-template/internal/repository"
	"simple-template/internal/utils"
	"simple-template/pkg/pagination"
	"strconv"
	"strings"
	"time"
)

type ProductUsecase struct {
	productRepo       *repository.ProductRepository
	paginationService *pagination.Service
}

func NewProductUsecase(productRepo *repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo:       productRepo,
		paginationService: pagination.NewService(),
	}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {

	if err := u.validateCreateProduct(req); err != nil {
		return nil, err
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
			Name:         strings.TrimSpace(reqVariant.Name),
			DisplayName:  strings.TrimSpace(reqVariant.DisplayName),
			IsRequire:    utils.DerefInt16OrDefault(reqVariant.IsRequire, 0),
			DisplayOrder: utils.DerefInt64OrDefault(reqVariant.DisplayOrder, 0),
			Price: model.ProductVariantPrice{
				Price:         utils.DerefFloat64OrDefault(reqVariant.Price, 0),
				Status:        1,
				EffectiveFrom: time.Now(),
			},
		}

		for _, reqValue := range reqVariant.Values {
			value := &model.ProductVariantValue{
				Value:         strings.TrimSpace(reqValue.Value),
				DisplayOrder:  utils.DerefIntOrDefault(reqValue.DisplayOrder, 0),
				StockQuantity: utils.DerefIntOrDefault(reqValue.StockQuantity, 0),
			}
			variant.Values = append(variant.Values, *value)
		}
		product.Variant = append(product.Variant, *variant)
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product %w", err)
	}

	return u.GetByID(ctx, product.ID)
}
func (u *ProductUsecase) validateCreateProduct(req *model.CreateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("product name is required")
	}

	if strings.TrimSpace(req.SKU) == "" {
		return fmt.Errorf("product SKU is required")
	}

	if req.CategoryID <= 0 {
		return fmt.Errorf("valid category ID is required")
	}

	if len(req.Variants) == 0 {
		return fmt.Errorf("at least one variant is required")
	}

	for i, variant := range req.Variants {
		if strings.TrimSpace(variant.Name) == "" {
			return fmt.Errorf("variant %d: name is required", i+1)
		}

		if variant.Price == nil || *variant.Price <= 0 {
			return fmt.Errorf("variant %d: valid price is required", i+1)
		}
	}

	return nil
}

func (u *ProductUsecase) GetByID(ctx context.Context, id int64) (*model.Product, error) {

	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (u *ProductUsecase) GetAll(ctx context.Context, req *pagination.Request) (*pagination.Response, error) {

	u.paginationService.ValidateAndNormalize(req)

	// Determine navigation direction (business logic)
	cursor, effectiveOrder := u.paginationService.GetNavigationParams(*req)

	// Fetch products (limit + 1 to check for next/prev page)
	fetchLimit := u.paginationService.CalculateFetchLimit(req.Limit)

	products, err := u.productRepo.GetAllPaginated(ctx, cursor, fetchLimit, effectiveOrder, req.SortBy)
	if err != nil {
		return nil, err
	}

	// variants
	var productIds []string
	for _, p := range products {
		productIds = append(productIds, strconv.FormatInt(p.ID, 10))
	}
	variants, err := u.productRepo.GetVariantsByProductIDs(ctx, productIds)
	if err != nil {
		return nil, err
	}
	// / variant values
	var variantIDs []string
	for _, v := range variants {
		variantIDs = append(variantIDs, strconv.FormatInt(v.ID, 10))
	}
	variantValues, err := u.productRepo.GetVariantValuesByAttributeID(ctx, variantIDs)
	if err != nil {
		return nil, err
	}

	// price
	prices, err := u.productRepo.GetPriceByVariantID(ctx, variantIDs)
	if err != nil {
		return nil, err
	}

	// combine data
	pricesMap := utils.ConvertArrToMapIDSlice(prices, func(v *model.Price) int64 { return v.VariantID })

	variantValueMap := utils.ConvertArrToMapIDSlice(variantValues, func(v *model.ProductVariantValue) int64 { return v.AttributeID })
	for _, v := range variants {
		v.Values = append(v.Values, variantValueMap[v.ID]...)
		if priceList, exists := pricesMap[v.ID]; exists && len(priceList) > 0 {
			v.Price.ID = priceList[0].ID
			v.Price.Price = priceList[0].Price
			v.Price.EffectiveFrom = priceList[0].EffectiveFrom
			v.Price.Status = int(priceList[0].Status)
		}
	}
	variantsMap := utils.ConvertArrToMapIDSlice(variants, func(v *model.ProductVariant) int64 { return v.ProductID })
	// Convert []*model.Product to []interface{} for pagination service
	items := make([]interface{}, len(products))

	for i, p := range products {
		p.Variant = append(p.Variant, variantsMap[p.ID]...)
		items[i] = p
	}

	response := u.paginationService.BuildResponse(
		items,
		req,
		func(p interface{}) (time.Time, int64) {
			product := p.(*model.Product)
			return product.CreatedAt, product.ID
		},
	)

	return &response, nil
}

func (u *ProductUsecase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid product id")
	}

	// Delete in correct order
	if err := u.productRepo.DeleteVariantValueByProductID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete variant values: %w", err)
	}

	if err := u.productRepo.DeletePriceByProductID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete prices: %w", err)
	}

	if err := u.productRepo.DeleteVariantByProductID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete variants: %w", err)
	}

	if err := u.productRepo.DeleteProductByID(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
