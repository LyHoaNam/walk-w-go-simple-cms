package repository

import (
	"context"
	"database/sql"
	"fmt"
	"simple-template/internal/database"
	"simple-template/internal/model"
	"simple-template/internal/utils"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type ProductRepository struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {

	// insert product
	queryProduct, arg, err := r.db.Dialect.Insert("product").Rows(goqu.Record{
		"name":        product.Name,
		"sku":         product.SKU,
		"status":      product.Status,
		"img_url":     product.ImgUrl,
		"category_id": product.CategoryID,
		"description": product.Description,
		"dimension":   product.Dimension,
		"weight":      product.Weight,
		"brand":       product.Brand,
		"material":    product.Material,
		"origin":      product.Origin,
	}).ToSQL()

	if err != nil {
		return fmt.Errorf("fail: %w", err)
	}
	result, err := r.db.SQL.ExecContext(ctx, queryProduct, arg...)
	if err != nil {
		return fmt.Errorf("fail: %w", err)
	}
	productID, err := result.LastInsertId()
	product.ID = productID
	if err != nil {
		return fmt.Errorf("fail: %w", err)
	}
	// insert variants
	for i, variant := range product.Variant {
		queryVariant, arg, err := r.db.Dialect.Insert("product_variant").Rows(goqu.Record{
			"name":          variant.Name,
			"display_name":  variant.DisplayName,
			"display_order": variant.DisplayOrder,
			"is_required":   variant.IsRequire,
			"product_id":    productID,
		}).ToSQL()

		if err != nil {
			return fmt.Errorf("fail: %w", err)
		}
		result, err = r.db.SQL.ExecContext(ctx, queryVariant, arg...)
		if err != nil {
			return fmt.Errorf("fail: %w", err)
		}
		variantID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("fail: %w", err)
		}
		product.Variant[i].ID = variantID

		for _, value := range variant.Values {
			queryValue, arg, err := r.db.Dialect.Insert("product_variant_value").Rows(goqu.Record{
				"attribute_id":   variantID,
				"value":          value.Value,
				"display_order":  value.DisplayOrder,
				"stock_quantity": value.StockQuantity,
			}).ToSQL()
			if err != nil {
				return fmt.Errorf("fail: %w", err)
			}
			result, err = r.db.SQL.ExecContext(ctx, queryValue, arg...)
			if err != nil {
				return fmt.Errorf("fail: %w", err)
			}
			valueID, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("fail: %w", err)
			}
			product.Variant[i].Values[i].ID = valueID
		}
	}

	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query, args, err := r.db.Dialect.Select(
		// Product fields
		goqu.I("product.id").As("product_id"),
		goqu.I("product.name").As("product_name"),
		goqu.I("product.sku"),
		goqu.I("product.status"),
		goqu.I("product.description"),
		goqu.I("product.dimension"),
		goqu.I("product.weight"),
		goqu.I("product.brand"),
		goqu.I("product.material"),
		goqu.I("product.origin"),
		goqu.I("product.img_url"),
		goqu.I("product.category_id"),
		goqu.I("product.created_at").As("product_created_at"),
		goqu.I("product.updated_at").As("product_updated_at"),
		// Product variant fields
		goqu.I("product_variant.id").As("variant_id"),
		goqu.I("product_variant.name").As("variant_name"),
		goqu.I("product_variant.display_name"),
		goqu.I("product_variant.display_order").As("variant_display_order"),
		goqu.I("product_variant.is_required"),
		goqu.I("product_variant.product_id"),
		// Product variant value fields
		goqu.I("product_variant_value.id").As("value_id"),
		goqu.I("product_variant_value.attribute_id"),
		goqu.I("product_variant_value.value"),
		goqu.I("product_variant_value.display_order").As("value_display_order"),
		goqu.I("product_variant_value.stock_quantity"),
		goqu.I("product_variant_value.created_at").As("value_created_at"),
		goqu.I("product_variant_value.updated_at").As("value_updated_at"),
	).
		From("product").
		LeftJoin(
			goqu.T("product_variant"),
			goqu.On(goqu.Ex{"product.id": goqu.I("product_variant.product_id")}),
		).
		LeftJoin(goqu.T("product_variant_value"),
			goqu.On(goqu.Ex{"product_variant.id": goqu.I("product_variant_value.attribute_id")})).
		Where(goqu.Ex{"product.id": id}).
		Order(goqu.I("product_variant.display_order").Asc(), goqu.I("product_variant_value.display_order").Asc()).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	defer rows.Close()

	var product *model.Product
	variantMap := make(map[int64]*model.ProductVariant)

	for rows.Next() {
		var (
			// Product fields
			productID        int64
			productName      string
			sku              string
			status           int32
			description      *string
			dimension        *string
			weight           *float64
			brand            *string
			material         *string
			origin           *string
			imgUrl           string
			categoryID       int64
			productCreatedAt time.Time
			productUpdatedAt time.Time
			// Variant fields (nullable)
			variantID        sql.NullInt64
			variantName      sql.NullString
			displayName      sql.NullString
			variantDispOrder sql.NullInt64
			isRequired       sql.NullInt16
			variantProductID sql.NullInt64
			// Variant value fields (nullable)
			valueID        sql.NullInt64
			attributeID    sql.NullInt64
			value          sql.NullString
			valueDispOrder sql.NullInt32
			stockQuantity  sql.NullInt32
			valueCreatedAt sql.NullTime
			valueUpdatedAt sql.NullTime
		)

		err := rows.Scan(
			&productID, &productName, &sku, &status, &description, &dimension,
			&weight, &brand, &material, &origin, &imgUrl, &categoryID,
			&productCreatedAt, &productUpdatedAt,
			&variantID, &variantName, &displayName, &variantDispOrder, &isRequired, &variantProductID,
			&valueID, &attributeID, &value, &valueDispOrder, &stockQuantity,
			&valueCreatedAt, &valueUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}

		// Initialize product on first row
		if product == nil {
			product = &model.Product{
				ID:          productID,
				Name:        productName,
				SKU:         sku,
				Status:      status,
				Description: description,
				Dimension:   dimension,
				Weight:      weight,
				Brand:       brand,
				Material:    material,
				Origin:      origin,
				ImgUrl:      imgUrl,
				CategoryID:  categoryID,
				CreatedAt:   productCreatedAt,
				UpdatedAt:   productUpdatedAt,
				Variant:     []model.ProductVariant{},
			}
		}

		// Add variant if exists
		if variantID.Valid {
			variant, exists := variantMap[variantID.Int64]
			if !exists {
				variant = &model.ProductVariant{
					ID:           variantID.Int64,
					Name:         variantName.String,
					DisplayName:  displayName.String,
					DisplayOrder: variantDispOrder.Int64,
					IsRequire:    isRequired.Int16,
					ProductID:    variantProductID.Int64,
					Values:       []model.ProductVariantValue{},
				}
				variantMap[variantID.Int64] = variant
				product.Variant = append(product.Variant, *variant)
			}

			// Add variant value if exists
			if valueID.Valid {
				variantValue := model.ProductVariantValue{
					ID:            valueID.Int64,
					AttributeID:   attributeID.Int64,
					Value:         value.String,
					DisplayOrder:  int(valueDispOrder.Int32),
					StockQuantity: int(stockQuantity.Int32),
					CreatedAt:     valueCreatedAt.Time,
					UpdatedAt:     valueUpdatedAt.Time,
				}
				variant.Values = append(variant.Values, variantValue)

				// Update the variant in the product slice
				for i := range product.Variant {
					if product.Variant[i].ID == variantID.Int64 {
						product.Variant[i] = *variant
						break
					}
				}
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}

func (r *ProductRepository) GetVariantValuesByAttributeID(ctx context.Context, attributeIDs []string) ([]*model.ProductVariantValue, error) {
	query, args, err := r.db.Dialect.
		Select("id", "attribute_id", "display_order", "stock_quantity", "value", "created_at", "updated_at").
		From("product_variant_value").
		Where(goqu.Ex{"attribute_id": attributeIDs}).
		Order(goqu.I("created_at").Desc()).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("fail to get variant attribute query: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("fail to select variant attribute value: %w", err)
	}
	var values []*model.ProductVariantValue
	for rows.Next() {
		var value model.ProductVariantValue
		err := rows.Scan(
			&value.ID,
			&value.AttributeID,
			&value.DisplayOrder,
			&value.StockQuantity,
			&value.Value,
			&value.CreatedAt,
			&value.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan fail variant value: %w", err)
		}
		values = append(values, &value)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return values, nil
}

func (r *ProductRepository) GetVariantsByProductIDs(ctx context.Context, productIDs []string) ([]*model.ProductVariant, error) {
	query, args, err := r.db.Dialect.
		Select("id", "name", "display_name", "display_order", "is_required", "product_id", "created_at", "updated_at").
		From("product_variant").
		Where(goqu.Ex{"product_id": productIDs}).
		Order(goqu.I("created_at").Desc()).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}
	var variants []*model.ProductVariant
	for rows.Next() {
		var variant model.ProductVariant
		err := rows.Scan(
			&variant.ID,
			&variant.Name,
			&variant.DisplayName,
			&variant.DisplayOrder,
			&variant.IsRequire,
			&variant.ProductID,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product variants: %w", err)
		}
		variants = append(variants, &variant)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return variants, nil
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]*model.Product, error) {
	query, args, err := r.db.Dialect.
		Select("id",
			"name",
			"sku",
			"description",
			"category_id",
			"created_at",
			"updated_at").From("product").Order(goqu.I("created_at").Desc()).ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.CategoryID,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	// variants
	var productIds []string
	for _, p := range products {
		productIds = append(productIds, strconv.FormatInt(p.ID, 10))
	}
	variants, err := r.GetVariantsByProductIDs(ctx, productIds)
	if err != nil {
		return nil, fmt.Errorf("failed to scan variants: %w", err)
	}
	// variant values
	var variantIDs []string
	for _, v := range variants {
		variantIDs = append(variantIDs, strconv.FormatInt(v.ID, 10))
	}
	variantValues, err := r.GetVariantValuesByAttributeID(ctx, variantIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to scan variants value: %w", err)

	}
	// combine data
	variantValueMap := utils.ConvertArrToMapIDSlice(variantValues, func(v *model.ProductVariantValue) int64 { return v.AttributeID })
	for _, v := range variants {
		v.Values = append(v.Values, variantValueMap[v.ID]...)
	}
	variantsMap := utils.ConvertArrToMapIDSlice(variants, func(v *model.ProductVariant) int64 { return v.ProductID })

	for _, p := range products {
		p.Variant = append(p.Variant, variantsMap[p.ID]...)
	}

	return products, nil
}

func (r *ProductRepository) DeleteByID(ctx context.Context, id int64) error {

	query, args, err := r.db.Dialect.Delete("product").Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build a delete query: %w", err)
	}
	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	// check the rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
