package usecase

import (
	"context"
	"fmt"
	"simple-template/internal/model"
	"simple-template/internal/repository"
	"strings"
)

type CustomerUsecase struct {
	customerRepo *repository.CustomerRepository
}

func NewCustomerUsecase(customerRepo *repository.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{
		customerRepo: customerRepo,
	}
}

func (u *CustomerUsecase) CreateCustomer(ctx context.Context, req *model.CreateCustomerRequest) (*model.Customer, error) {
	err := u.validateCreateCustomer(req)
	if err != nil {
		return nil, err
	}

	customer := &model.Customer{
		FirstName:   strings.TrimSpace(req.FirstName),
		LastName:    strings.TrimSpace(req.LastName),
		Address:     strings.TrimSpace(req.Address),
		Email:       strings.ToLower(strings.TrimSpace(req.Email)),
		PhoneNumber: req.PhoneNumber,
	}

	result, err := u.customerRepo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *CustomerUsecase) validateCreateCustomer(req *model.CreateCustomerRequest) error {
	if strings.TrimSpace(req.FirstName) == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.PhoneNumber) < 10 || len(req.PhoneNumber) > 11 {
		return fmt.Errorf("invalid phone number")
	}
	return nil
}

func (u *CustomerUsecase) validateUpdateCustomer(req *model.UpdateCustomerRequest) error {
	if req.PhoneNumber != nil && (len(*req.PhoneNumber) < 10 || len(*req.PhoneNumber) > 11) {
		return fmt.Errorf("invalid phone number")
	}
	return nil
}

func (u *CustomerUsecase) GetCustomerByID(ctx context.Context, id int64) (*model.Customer, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}
	customer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (u *CustomerUsecase) GetAllCustomer(ctx context.Context) ([]*model.Customer, error) {
	customer, err := u.customerRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (u *CustomerUsecase) UpdateCustomer(ctx context.Context, id int64, req *model.UpdateCustomerRequest) (*model.Customer, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}
	if err := u.validateUpdateCustomer(req); err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Address != nil && strings.TrimSpace(*req.Address) != "" {
		updates["address"] = strings.TrimSpace(*req.Address)
	}
	if req.Email != nil && strings.TrimSpace(*req.Email) != "" {
		updates["email"] = strings.TrimSpace(*req.Email)
	}
	if req.FirstName != nil && strings.TrimSpace(*req.FirstName) != "" {
		updates["first_name"] = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil && strings.TrimSpace(*req.LastName) != "" {
		updates["last_name"] = strings.TrimSpace(*req.LastName)
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = strings.TrimSpace(*req.PhoneNumber)
	}

	if err := u.customerRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}

	updatedCustomer, err := u.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedCustomer, nil
}

func (u *CustomerUsecase) DeleteCustomer(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid id")
	}

	if err := u.customerRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
