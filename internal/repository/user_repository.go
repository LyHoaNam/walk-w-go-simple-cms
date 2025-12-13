package repository

import (
	"context"
	"database/sql"
	"fmt"

	"simple-template/internal/database"
	"simple-template/internal/model"

	"github.com/doug-martin/goqu/v9"
)

// UserRepository handles all database operations related to users
type UserRepository struct {
	db *database.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	// Build insert query with goqu
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

	// Get the ID of the created user
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	return nil
}

// GetByID gets user information by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	// Build select query with goqu
	query, args, err := r.db.Dialect.
		Select("id", "name", "email", "created_at", "updated_at").
		From("users").
		Where(goqu.Ex{"id": id}).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	// Execute query
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

// GetAll gets a list of all users
func (r *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	// Build select query with goqu
	query, args, err := r.db.Dialect.
		Select("id", "name", "email", "created_at", "updated_at").
		From("users").
		Order(goqu.I("created_at").Desc()).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %w", err)
	}

	// Execute query
	rows, err := r.db.SQL.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	// Read results
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

// Update updates user information
func (r *UserRepository) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	// Build update query with goqu
	query, args, err := r.db.Dialect.
		Update("users").
		Set(updates).
		Where(goqu.Ex{"id": id}).
		ToSQL()

	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := r.db.SQL.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	// Build delete query with goqu
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
