package repository

import (
	"context"
	"database/sql"
	"fmt"
	"simple-template/internal/database"
	"simple-template/internal/model"

	"github.com/doug-martin/goqu/v9"
)

type OrdersRepository struct {
	db *database.DB
}

func NewOrdersRepository(db *database.DB) *OrdersRepository {
	return &OrdersRepository{
		db: db,
	}
}

func (r *OrdersRepository) Create(ctx context.Context, tx *sql.Tx, orders *model.Orders) (*model.Orders, error) {
	query, args, err := r.db.Dialect.
		Insert("orders").Rows(
		goqu.Record{
			"payment_status":   orders.PaymentStatus,
			"customer_id":      orders.CustomerID,
			"platform_id":      orders.PlatformID,
			"payment_id":       orders.PaymentID,
			"retail_stores_id": orders.RetailStoreID,
		}).ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	orders.ID = id

	return orders, nil
}

func (r *OrdersRepository) CreateItems(ctx context.Context, tx *sql.Tx, orderItems []*model.OrderItems) ([]*model.OrderItems, error) {
	itemRecords := make([]interface{}, 0, len(orderItems))
	for _, item := range orderItems {
		itemRecords = append(itemRecords, goqu.Record{
			"order_id":         item.OrderID,
			"price_id":         item.PriceID,
			"quantity":         item.Quantity,
			"variant_value_id": item.VariantValueID,
		})
	}

	query, args, err := r.db.Dialect.
		Insert("order_items").Rows(itemRecords...).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build order items query: %w", err)
	}
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order item: %w", err)
	}

	firstItemID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last order item: %w", err)
	}
	// assign ID
	for i := range orderItems {
		orderItems[i].ID = firstItemID + int64(i)
	}
	return orderItems, nil
}

func (r *OrdersRepository) GetStocks(
	ctx context.Context,
	PricesIds []int64,
	variantValueIds []int64,
) ([]*model.OrdersProduct, error) {
	query, args, err := r.db.Dialect.
		Select(
			goqu.I("price.ID").As("price_id"),
			goqu.I("product_variant.ID").As("product_variant_id"),
			goqu.I("product_variant_value.ID").As("product_variant_value_id"),
			goqu.I("product_variant_value.stock_quantity").As("stock_quantity"),
			goqu.I("product_variant_value.value").As("value"),
			goqu.I("product_variant.name").As("name"),
			goqu.I("price.status").As("status"),
		).From("price").
		LeftJoin(
			goqu.T("product_variant"),
			goqu.On(goqu.Ex{"price.variant_id": goqu.I("product_variant.id")}),
		).LeftJoin(
		goqu.T("product_variant_value"),
		goqu.On(goqu.Ex{"product_variant.id": goqu.I("product_variant_value.attribute_id")}),
	).Where(
		goqu.Ex{
			"price.id":                 PricesIds,
			"product_variant_value.id": variantValueIds,
		}).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build the query: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var OrdersProducts []*model.OrdersProduct

	for rows.Next() {
		OP := &model.OrdersProduct{}
		err := rows.Scan(
			&OP.PriceID,
			&OP.VariantID,
			&OP.VariantValueID,
			&OP.StockQuantity,
			&OP.Value,
			&OP.Name,
			&OP.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		OrdersProducts = append(OrdersProducts, OP)
	}

	return OrdersProducts, nil
}

func (r *OrdersRepository) ReduceStocksBatch(
	ctx context.Context,
	tx *sql.Tx,
	stockUpdates map[int64]int64,
) error {
	if len(stockUpdates) == 0 {
		return nil
	}

	for variantValueID, quantityToReduce := range stockUpdates {
		query, args, err := r.db.Dialect.
			Update("product_variant_value").
			Set(goqu.Record{
				"stock_quantity": goqu.L("stock_quantity - ?", quantityToReduce),
			}).
			Where(goqu.Ex{"id": variantValueID}).
			ToSQL()

		if err != nil {
			return fmt.Errorf("failed to build query: %w", err)
		}

		result, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		rowAffected, err := result.RowsAffected()
		if err != nil || rowAffected == 0 {
			return fmt.Errorf("failed to update")
		}
	}

	return nil
}

// BeginTx starts a new transaction
func (r *OrdersRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.SQL.BeginTx(ctx, nil)
}

func (r *OrdersRepository) GetOrdersPage(ctx context.Context) ([]*model.OrdersPage, error) {
	// Subquery: get latest order status
	latestStatusSubquery := r.db.Dialect.
		Select(
			goqu.I("order_id"),
			goqu.I("status"),
			goqu.I("description"),
			goqu.L("ROW_NUMBER() OVER (PARTITION BY order_id ORDER BY created_at DESC)").As("rn"),
		).
		From("order_status")

	// Subquery: calculate order totals
	orderTotalsSubquery := r.db.Dialect.
		Select(
			goqu.I("oi.order_id"),
			goqu.L("SUM(oi.quantity * p.price)").As("total_amount"),
		).From(goqu.T("order_items").As("oi")).
		LeftJoin(
			goqu.T("price").As("p"),
			goqu.On(goqu.Ex{"p.id": goqu.I("oi.price_id")}),
		).
		GroupBy("oi.order_id")

	query, args, err := r.db.Dialect.
		Select(
			goqu.I("orders.id"),
			goqu.I("orders.payment_status"),
			goqu.I("orders.created_at"),
			goqu.I("customer.first_name"),
			goqu.I("customer.last_name"),
			goqu.I("platform.name").As("platform"),
			goqu.I("payment_methods.name").As("payment_method"),
			goqu.I("ls.status").As("order_status"),
			goqu.L("COALESCE(ot.total_amount, 0)").As("total_amount"),
		).
		From("orders").
		LeftJoin(
			goqu.T("customer"),
			goqu.On(goqu.Ex{"orders.customer_id": goqu.I("customer.id")}),
		).
		LeftJoin(
			goqu.T("payment_methods"),
			goqu.On(goqu.Ex{"orders.payment_id": goqu.I("payment_methods.id")}),
		).
		LeftJoin(
			goqu.T("platform"),
			goqu.On(goqu.Ex{"orders.platform_id": goqu.I("platform.id")}),
		).
		LeftJoin(
			latestStatusSubquery.As("ls"),
			goqu.On(goqu.And(
				goqu.Ex{"ls.order_id": goqu.I("orders.id")},
				goqu.Ex{"ls.rn": 1},
			)),
		).
		LeftJoin(
			orderTotalsSubquery.As("ot"),
			goqu.On(goqu.Ex{"ot.order_id": goqu.I("orders.id")})).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var orders []*model.OrdersPage
	for rows.Next() {
		order := &model.OrdersPage{}
		err := rows.Scan(
			&order.ID,
			&order.PaymentStatus,
			&order.CreatedAt,
			&order.FirstName,
			&order.LastName,
			&order.Platform,
			&order.PaymentMethod,
			&order.OrderStatus,
			&order.TotalAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return orders, nil
}
