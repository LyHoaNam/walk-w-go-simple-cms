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
	if err := u.orderRepo.CreateOrderStatus(ctx, tx, &model.OrderStatus{
		Status:      1,
		Description: "created new orders",
		OrderID:     orders.ID,
	}); err != nil {
		return nil, fmt.Errorf("failed to create order status: %w", err)
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

func (u *OrderUsecase) GetOrdersPage(ctx context.Context) ([]*model.OrdersPage, error) {
	orders, err := u.orderRepo.GetOrdersPage(ctx)

	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (u *OrderUsecase) UpdateOrderStatus(
	ctx context.Context,
	status int8,
	orderID int64) error {

	if status < 1 || status > 5 {
		return fmt.Errorf("invalid status: must be between 1-5")
	}

	// Start transaction for status update
	tx, err := u.orderRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	latestStatus, err := u.orderRepo.GetLatestStatus(ctx, orderID)
	if err != nil {
		return err
	}

	// Validate status transition
	if err := u.validateStatusTransition(latestStatus.Status, status); err != nil {
		return err
	}

	// Create new order status record
	newStatus := &model.OrderStatus{
		Status:      status,
		Description: u.getStatusDescription(status),
		OrderID:     orderID,
	}

	if err := u.orderRepo.CreateOrderStatus(ctx, tx, newStatus); err != nil {
		return fmt.Errorf("failed to create order status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// validateStatusTransition checks if status transition is allowed
func (u *OrderUsecase) validateStatusTransition(currentStatus, newStatus int8) error {
	// Can't transition to the same status
	if currentStatus == newStatus {
		return fmt.Errorf("order is already in status %d", currentStatus)
	}

	// Define valid transitions map
	validTransitions := map[int8][]int8{
		int8(model.OrderStatusPending): {
			int8(model.OrderStatusPaid),
			int8(model.OrderStatusCanceled),
		},
		int8(model.OrderStatusPaid): {
			int8(model.OrderStatusShipped),
			int8(model.OrderStatusCanceled),
		},
		int8(model.OrderStatusShipped): {
			int8(model.OrderStatusCompleted),
			int8(model.OrderStatusCanceled),
		},
		int8(model.OrderStatusCompleted): {}, // Final state - no transitions allowed
		int8(model.OrderStatusCanceled):  {}, // Final state - no transitions allowed
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("unknown current status: %d", currentStatus)
	}

	// Check if new status is in allowed transitions
	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return nil // Valid transition
		}
	}

	// Invalid transition
	return fmt.Errorf("cannot transition from status %d to %d", currentStatus, newStatus)
}

// getStatusDescription returns description based on status
func (u *OrderUsecase) getStatusDescription(status int8) string {
	switch status {
	case int8(model.OrderStatusPending):
		return "Order is pending payment"
	case int8(model.OrderStatusPaid):
		return "Payment confirmed successfully"
	case int8(model.OrderStatusShipped):
		return "Order has been shipped"
	case int8(model.OrderStatusCompleted):
		return "Order completed and delivered"
	case int8(model.OrderStatusCanceled):
		return "Order has been canceled"
	default:
		return "Status updated"
	}
}
