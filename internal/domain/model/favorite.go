package model

type Favorite struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	ProductID  string `json:"product_id"`
	CreatedAt  string `json:"created_at"`
}
