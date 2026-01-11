package pagination

import (
	"github.com/doug-martin/goqu/v9"
)

// QueryBuilder builds cursor-based pagination queries for goqu
type QueryBuilder struct{}

// NewQueryBuilder creates a new query builder instance
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// ApplyCursorPagination applies sorting and cursor filtering to a goqu query
// This method is generic and works for any table with id and a timestamp field
//
// Parameters:
//   - query: The base goqu SelectDataset to apply pagination to
//   - cursor: Encoded cursor string (will be decoded internally)
//   - limit: Number of items to fetch
//   - order: Sort direction ("asc" or "desc")
//   - sortBy: Field to sort by ("id" or a timestamp field like "created_at")
//
// Returns: Modified query with ORDER BY and WHERE clauses applied
func (qb *QueryBuilder) ApplyCursorPagination(
	query *goqu.SelectDataset,
	cursor string,
	limit int,
	order string,
	sortBy string,
) (*goqu.SelectDataset, error) {
	// Decode cursor
	cursorTimestamp, cursorID, err := DecodeCursor(cursor)
	if err != nil {
		return nil, err
	}

	// Apply limit
	query = query.Limit(uint(limit))

	// Apply sorting and filtering based on sortBy field
	if sortBy == "id" {
		return qb.applySortByID(query, order, cursorID), nil
	}

	// Default: sort by timestamp field (e.g., created_at, updated_at)
	return qb.applySortByTimestamp(query, order, sortBy, cursorTimestamp, cursorID), nil
}

// applySortByID handles simple ID-based sorting and filtering
func (qb *QueryBuilder) applySortByID(
	query *goqu.SelectDataset,
	order string,
	cursorID int64,
) *goqu.SelectDataset {
	if order == "asc" {
		query = query.Order(goqu.I("id").Asc())
		if cursorID > 0 {
			query = query.Where(goqu.I("id").Gt(cursorID))
		}
	} else {
		query = query.Order(goqu.I("id").Desc())
		if cursorID > 0 {
			query = query.Where(goqu.I("id").Lt(cursorID))
		}
	}
	return query
}

// applySortByTimestamp handles composite sorting by timestamp + ID
// This prevents missing/duplicate items when multiple records share the same timestamp
func (qb *QueryBuilder) applySortByTimestamp(
	query *goqu.SelectDataset,
	order string,
	sortByField string,
	cursorTimestamp string,
	cursorID int64,
) *goqu.SelectDataset {
	// Create field identifier for the sort column
	timestampField := goqu.I(sortByField)

	if order == "asc" {
		// Sort: oldest first, then by ID ascending
		query = query.Order(timestampField.Asc(), goqu.I("id").Asc())

		// Filter: get records after cursor
		// WHERE (timestamp > cursor_time) OR (timestamp = cursor_time AND id > cursor_id)
		if cursorTimestamp != "" {
			query = query.Where(
				goqu.Or(
					timestampField.Gt(cursorTimestamp),
					goqu.And(
						timestampField.Eq(cursorTimestamp),
						goqu.I("id").Gt(cursorID),
					),
				),
			)
		}
	} else {
		// Sort: newest first, then by ID descending
		query = query.Order(timestampField.Desc(), goqu.I("id").Desc())

		// Filter: get records before cursor
		// WHERE (timestamp < cursor_time) OR (timestamp = cursor_time AND id < cursor_id)
		if cursorTimestamp != "" {
			query = query.Where(
				goqu.Or(
					timestampField.Lt(cursorTimestamp),
					goqu.And(
						timestampField.Eq(cursorTimestamp),
						goqu.I("id").Lt(cursorID),
					),
				),
			)
		}
	}

	return query
}

// ApplyCursorPaginationWithTablePrefix is a variant that supports table prefixes
// Useful for queries with JOINs where column names need qualification (e.g., "product.id")
//
// Example: ApplyCursorPaginationWithTablePrefix(query, cursor, 10, "desc", "created_at", "product")
func (qb *QueryBuilder) ApplyCursorPaginationWithTablePrefix(
	query *goqu.SelectDataset,
	cursor string,
	limit int,
	order string,
	sortBy string,
	tablePrefix string,
) (*goqu.SelectDataset, error) {
	// Decode cursor
	cursorTimestamp, cursorID, err := DecodeCursor(cursor)
	if err != nil {
		return nil, err
	}

	// Apply limit
	query = query.Limit(uint(limit))

	// Build qualified field names
	idField := goqu.I(tablePrefix + ".id")
	timestampField := goqu.I(tablePrefix + "." + sortBy)

	if sortBy == "id" {
		if order == "asc" {
			query = query.Order(idField.Asc())
			if cursorID > 0 {
				query = query.Where(idField.Gt(cursorID))
			}
		} else {
			query = query.Order(idField.Desc())
			if cursorID > 0 {
				query = query.Where(idField.Lt(cursorID))
			}
		}
		return query, nil
	}

	// Sort by timestamp + id (composite)
	if order == "asc" {
		query = query.Order(timestampField.Asc(), idField.Asc())
		if cursorTimestamp != "" {
			query = query.Where(
				goqu.Or(
					timestampField.Gt(cursorTimestamp),
					goqu.And(
						timestampField.Eq(cursorTimestamp),
						idField.Gt(cursorID),
					),
				),
			)
		}
	} else {
		query = query.Order(timestampField.Desc(), idField.Desc())
		if cursorTimestamp != "" {
			query = query.Where(
				goqu.Or(
					timestampField.Lt(cursorTimestamp),
					goqu.And(
						timestampField.Eq(cursorTimestamp),
						idField.Lt(cursorID),
					),
				),
			)
		}
	}

	return query, nil
}
