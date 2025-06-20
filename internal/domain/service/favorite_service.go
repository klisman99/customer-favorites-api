package service

import (
	"app/internal/domain/model"
	"app/internal/infra/db"
	"context"
)

type FavoriteService struct {
	favoriteRepo   *db.FavoriteRepository
	productService ProductService
}

func NewFavoriteService(favoriteRepo *db.FavoriteRepository, productService ProductService) *FavoriteService {
	return &FavoriteService{
		favoriteRepo:   favoriteRepo,
		productService: productService,
	}
}
func (s *FavoriteService) AddFavorite(c context.Context, customerID string, productID int) error {
	return s.favoriteRepo.AddFavorite(c, customerID, productID)
}

func (s *FavoriteService) RemoveFavorite(c context.Context, productID int) error {
	return s.favoriteRepo.RemoveFavorite(c, productID)
}

func (s *FavoriteService) GetCustomerFavoriteProducts(c context.Context, customerID string) ([]model.Product, error) {
	productIDs, err := s.favoriteRepo.FindProductsByCustomerID(c, customerID)
	if err != nil {
		return nil, err
	}

	products, err := s.productService.GetAll()
	if err != nil {
		return nil, err
	}

	var favorites []model.Product
	for _, id := range productIDs {
		for _, p := range products {
			if p.ID == id {
				favorites = append(favorites, p)
			}
		}
	}

	return favorites, nil
}
