package service

import (
	"fmt"
	"strings"

	"github.com/emmrys-jay/gigmile/internal/models"
	"github.com/emmrys-jay/gigmile/internal/repository"
)

type CustomerService interface {
	CreateCustomer(customerReq *models.CreateCustomerRequest) (*models.Customer, error)
	GetCustomerByID(id int64) (*models.Customer, error)
	GetAllCustomers() ([]*models.Customer, error)
	UpdateCustomer(id int64, customerReq *models.UpdateCustomerRequest) (*models.Customer, error)
	DeleteCustomer(id int64) error
}

type customerService struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository) CustomerService {
	return &customerService{
		customerRepo: customerRepo,
	}
}

func (s *customerService) CreateCustomer(customerReq *models.CreateCustomerRequest) (*models.Customer, error) {
	// Normalize email
	customerReq.Email = strings.ToLower(strings.TrimSpace(customerReq.Email))

	// Create customer
	customer, err := s.customerRepo.Create(customerReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

func (s *customerService) GetCustomerByID(id int64) (*models.Customer, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid customer id")
	}

	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) GetAllCustomers() ([]*models.Customer, error) {
	customers, err := s.customerRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (s *customerService) UpdateCustomer(id int64, customerReq *models.UpdateCustomerRequest) (*models.Customer, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid customer id")
	}

	// Normalize email if provided
	if customerReq.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*customerReq.Email))
		customerReq.Email = &email
	}

	customer, err := s.customerRepo.Update(id, customerReq)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) DeleteCustomer(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid customer id")
	}

	err := s.customerRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
