package repository

import (
	"context"
	"database/sql"
	"fmt"

	"simple-template/internal/database"
	"simple-template/internal/model"

	"github.com/doug-martin/goqu/v9"
)

// UserRepository xử lý tất cả các thao tác database liên quan đến user
type UserRepository struct {
	db *database.DB
}

// NewUserRepository tạo instance mới của UserRepository
func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create tạo user mới trong database
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	// Tạo query insert với goqu
	query, args, err := r.db.Dialect.
		Insert("users").
		Rows(goqu.Record{
			"name":  user.Name,
			"email": user.Email,
		}).
		ToSQL()

	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	// Thực thi query
	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Lấy ID của user vừa tạo
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	return nil
}

// GetByID lấy thông tin user theo ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	// Tạo query select với goqu
	query, args, err := r.db.Dialect.
		Select("id", "name", "email", "created_at", "updated_at").
		From("users").
		Where(goqu.Ex{"id": id}).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	// Thực thi query
	var user model.User
	err = r.db.SQL.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetAll lấy danh sách tất cả users
func (r *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	// Tạo query select với goqu
	query, args, err := r.db.Dialect.
		Select("id", "name", "email", "created_at", "updated_at").
		From("users").
		Order(goqu.I("created_at").Desc()).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	// Thực thi query
	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	// Đọc kết quả
	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return users, nil
}

// Update cập nhật thông tin user
func (r *UserRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	// Tạo query update với goqu
	query, args, err := r.db.Dialect.
		Update("users").
		Set(updates).
		Where(goqu.Ex{"id": id}).
		ToSQL()

	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// Thực thi query
	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Kiểm tra số dòng bị ảnh hưởng
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete xóa user theo ID
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	// Tạo query delete với goqu
	query, args, err := r.db.Dialect.
		Delete("users").
		Where(goqu.Ex{"id": id}).
		ToSQL()

	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	// Thực thi query
	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Kiểm tra số dòng bị ảnh hưởng
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
