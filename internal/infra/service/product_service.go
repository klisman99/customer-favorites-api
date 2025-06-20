package service

import (
	"app/internal/domain/model"
	"encoding/json"
	"net/http"
)

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) GetAll() ([]model.Product, error) {
	resp, err := http.Get("https://fakestoreapi.com/products")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []model.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	return products, nil
}
