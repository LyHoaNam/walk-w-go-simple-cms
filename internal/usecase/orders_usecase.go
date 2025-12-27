package usecase

import (
	"context"
	"fmt"
	"simple-template/internal/model"
	"simple-template/internal/repository"
)

type OrderUsecase struct {
	orderRepo *repository.OrdersRepository
}

func NewOrderUseCase(orderRepo *repository.OrdersRepository) *OrderUsecase {
	return &OrderUsecase{
		orderRepo: orderRepo,
	}
}

func (u *OrderUsecase) CreateOrders(ctx context.Context, req *model.CreateOrders) (*model.Orders, error) {
	// Validate before starting transaction
	if err := u.validateCreateOrder(ctx, req); err != nil {
		return nil, err
	}

	// Start transaction
	tx, err := u.orderRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Safety net: rollback if function exits without commit
	defer tx.Rollback()

	// Create order within transaction
	createOrders := &model.Orders{
		PaymentStatus: int8(1),
		CustomerID:    req.CustomerID,
		PlatformID:    req.PlatformID,
		RetailStoreID: req.RetailStoreID,
		PaymentID:     req.PaymentID,
	}

	orders, err := u.orderRepo.Create(ctx, tx, createOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items within transaction
	var createItems []*model.OrderItems
	for _, item := range req.Items {
		createItems = append(createItems, &model.OrderItems{
			OrderID:        orders.ID,
			PriceID:        item.PriceID,
			VariantValueID: item.ProductVariantID,
			Quantity:       int(item.Quantity),
		})
	}

	items, err := u.orderRepo.CreateItems(ctx, tx, createItems)
	if err != nil {
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}

	// Reduce stock within transaction (using VariantValueID, not item.ID)
	stockUpdates := make(map[int64]int64)
	for _, item := range req.Items {
		stockUpdates[item.ProductVariantID] = item.Quantity
	}

	if err := u.orderRepo.ReduceStocksBatch(ctx, tx, stockUpdates); err != nil {
		return nil, fmt.Errorf("failed to reduce stocks: %w", err)
	}

	// Commit transaction - all operations succeeded atomically
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Combine results
	orders.Items = items

	return orders, nil
}

func (u *OrderUsecase) validateCreateOrder(ctx context.Context, req *model.CreateOrders) error {
	if len(req.Items) <= 0 {
		return fmt.Errorf("invalid request")
	}
	var (
		ProductVariantIDs []int64
		priceIDs          []int64
	)
	stockQuantityMap := make(map[int64]int64)

	for _, r := range req.Items {
		ProductVariantIDs = append(ProductVariantIDs, r.ProductVariantID)
		priceIDs = append(priceIDs, r.PriceID)
		stockQuantityMap[r.ProductVariantID] = r.Quantity
	}

	stocks, err := u.orderRepo.GetStocks(ctx, priceIDs, ProductVariantIDs)

	if err != nil {
		return err
	}
	if len(stocks) == 0 {
		return fmt.Errorf("price don't match with variant")
	}

	for _, stock := range stocks {
		if stock.Status != 1 {
			return fmt.Errorf("the product %s, %v have status inactive", stock.Name, stock.Value)
		}
		requestQuantity, exist := stockQuantityMap[stock.VariantValueID]
		if !exist {
			return fmt.Errorf("price don't match the variant value")
		}
		if requestQuantity > int64(stock.StockQuantity) {
			return fmt.Errorf("out of stock")
		}
	}
	return nil
}
