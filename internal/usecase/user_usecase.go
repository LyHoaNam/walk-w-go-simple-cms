package usecase

import (
	"context"
	"fmt"
	"strings"

	"simple-template/internal/model"
	"simple-template/internal/repository"
)

// UserUsecase handles business logic related to users
type UserUsecase struct {
	userRepo *repository.UserRepository
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (u *UserUsecase) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Validate input
	if err := u.validateCreateUser(req); err != nil {
		return nil, err
	}

	// Create user object
	user := &model.User{
		Name:  strings.TrimSpace(req.Name),
		Email: strings.ToLower(strings.TrimSpace(req.Email)),
	}

	// Call repository to save to database
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get the created user info (to have created_at, updated_at)
	createdUser, err := u.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created user: %w", err)
	}

	return createdUser, nil
}

// GetUserByID gets user information by ID
func (u *UserUsecase) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	// Call repository to get user
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers gets a list of all users
func (u *UserUsecase) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	// Call repository to get list of users
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates user information
func (u *UserUsecase) UpdateUser(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	// Validate input
	if err := u.validateUpdateUser(req); err != nil {
		return nil, err
	}

	// Tạo map updates
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}
	if req.Email != "" {
		updates["email"] = strings.ToLower(strings.TrimSpace(req.Email))
	}

	// Check if there's anything to update
	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Call repository to update
	if err := u.userRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}

	// Get user info after update
	updatedUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	return updatedUser, nil
}

// DeleteUser deletes a user by ID
func (u *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid user id")
	}

	// Call repository to delete user
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// validateCreateUser validate dữ liệu khi tạo user
func (u *UserUsecase) validateCreateUser(req *model.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}

	// Validate email format đơn giản
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// validateUpdateUser validate dữ liệu khi update user
func (u *UserUsecase) validateUpdateUser(req *model.UpdateUserRequest) error {
	// Nếu có email, validate format
	if req.Email != "" && !strings.Contains(req.Email, "@") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
