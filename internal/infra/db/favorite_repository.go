package db

import (
	"context"
	"database/sql"
)

type FavoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) *FavoriteRepository {
	return &FavoriteRepository{
		db: db,
	}
}

func (r *FavoriteRepository) AddFavorite(c context.Context, customerID string, productID int) error {
	query := `
		INSERT INTO customers_favorite_products (customer_id, product_id)
		VALUES ($1, $2)
		ON CONFLICT (customer_id, product_id) DO NOTHING
	`
	_, err := r.db.ExecContext(c, query, customerID, productID)
	if err != nil {
		return err
	}
	return nil
}

func (r *FavoriteRepository) RemoveFavorite(c context.Context, productID int) error {
	query := `
		DELETE FROM customers_favorite_products
		WHERE product_id = $1
	`
	_, err := r.db.ExecContext(c, query, productID)
	if err != nil {
		return err
	}
	return nil
}

func (r *FavoriteRepository) FindProductsByCustomerID(c context.Context, customerID string) ([]int, error) {
	query := `
		SELECT product_id
		FROM customers_favorite_products
		WHERE customer_id = $1
	`
	rows, err := r.db.QueryContext(c, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productIDs []int
	for rows.Next() {
		var productID int
		if err := rows.Scan(&productID); err != nil {
			return nil, err
		}
		productIDs = append(productIDs, productID)
	}

	return productIDs, nil
}
