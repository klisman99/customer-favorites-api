package service

import (
	"app/internal/domain"
	"app/internal/domain/model"
	"app/internal/infra/db"
	"context"

	"github.com/google/uuid"
)

type CustomerService struct {
	customerRepo *db.CustomerRepository
}

func NewCustomerService(repo *db.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: repo}
}

func (s *CustomerService) Create(c context.Context, name string, email string) (*model.Customer, error) {
	emailExists, err := s.customerRepo.FindByEmail(c, email, "")
	if err != nil {
		return nil, err
	}
	if emailExists != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	customer := model.Customer{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}

	createdCustomer, err := s.customerRepo.Create(c, customer)
	if err != nil {
		return nil, err
	}

	return createdCustomer, nil
}

func (s *CustomerService) GetByID(c context.Context, id string) (*model.Customer, error) {
	return s.customerRepo.FindByID(c, id)
}

func (s *CustomerService) GetAll(c context.Context) ([]model.Customer, error) {
	return s.customerRepo.FindAll(c)
}

func (s *CustomerService) Update(c context.Context, id string, name string, email string) error {
	emailExists, err := s.customerRepo.FindByEmail(c, email, id)
	if err != nil {
		return err
	}
	if emailExists != nil {
		return domain.ErrEmailAlreadyExists
	}

	customer, err := s.customerRepo.FindByID(c, id)
	if err != nil {
		return err
	}

	if customer == nil {
		return domain.ErrNotFound
	}

	customer.Name = name
	customer.Email = email

	return s.customerRepo.Update(c, *customer)
}

func (s *CustomerService) Delete(c context.Context, id string) error {
	customer, err := s.customerRepo.FindByID(c, id)
	if err != nil {
		return err
	}

	if customer == nil {
		return domain.ErrNotFound
	}

	return s.customerRepo.Delete(c, id)
}
