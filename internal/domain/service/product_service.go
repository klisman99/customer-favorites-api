package service

import (
	"app/internal/domain/model"
)

type ProductService interface {
	GetAll() ([]model.Product, error)
}
