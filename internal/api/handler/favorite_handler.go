package handler

import (
	"app/internal/domain/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FavoriteHandler struct {
	favoriteService *service.FavoriteService
}

type FavoriteIncludeRequest struct {
	ProductID int `json:"product_id" validate:"required" example:"123"`
}

type Product struct {
	ID          int     `json:"id" example:"123"`
	Name        string  `json:"name" example:"Smartphone XYZ"`
	Description string  `json:"description,omitempty" example:"Latest smartphone model"`
	Price       float64 `json:"price" example:"999.99"`
	Category    string  `json:"category,omitempty" example:"Electronics"`
}

func NewFavoriteHandler(favoriteService *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

// @Summary Add a product to customer's favorites
// @Description Adds a product to the specified customer's list of favorite products
// @Tags Favorite
// @Accept json
// @Produce json
// @Param customer_id path string true "Customer ID" example="550e8400-e29b-41d4-a716-446655440000"
// @Param favorite body FavoriteIncludeRequest true "Product to add to favorites"
// @Success 204 "Product added to favorites"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 404 {object} ErrorResponse "Customer not found"
// @Failure 500 {object} ErrorResponse "Failed to add to favorites"
// @Router /api/v1/customers/{customer_id}/favorites [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	customerID := c.Param("customer_id")

	var req FavoriteIncludeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding request data", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	if err := validator.New().Struct(req); err != nil {
		log.Println("Error validating request data", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	err := h.favoriteService.AddFavorite(c, customerID, req.ProductID)
	if err != nil {
		log.Println("Error adding favorite", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add favorite"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Remove a product from favorites
// @Description Removes a product from the customer's list of favorite products
// @Tags Favorite
// @Produce json
// @Param customer_id path string true "Customer ID" example="550e8400-e29b-41d4-a716-446655440000"
// @Param product_id path int true "Product ID to remove from favorites" example=123
// @Success 204 "Product removed from favorites"
// @Failure 400 {object} ErrorResponse "Invalid product ID"
// @Failure 404 {object} ErrorResponse "Product not found in favorites"
// @Failure 500 {object} ErrorResponse "Failed to remove from favorites"
// @Router /api/v1/customers/{customer_id}/favorites/{product_id} [delete]
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	productID := c.Param("product_id")

	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	intProductID, err := strconv.Atoi(productID)
	if err != nil {
		log.Println("Error converting product id to int", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	err = h.favoriteService.RemoveFavorite(c, intProductID)
	if err != nil {
		log.Println("Error removing favorite", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove favorite"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get customer's favorite products
// @Description Retrieves a list of all favorite products for the specified customer
// @Tags Favorite
// @Produce json
// @Param customer_id path string true "Customer ID" example="550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {array} Product "List of favorite products"
// @Failure 400 {object} ErrorResponse "Invalid customer ID"
// @Failure 404 {object} ErrorResponse "Customer not found"
// @Failure 500 {object} ErrorResponse "Failed to fetch favorite products"
// @Router /api/v1/customers/{customer_id}/favorites [get]
func (h *FavoriteHandler) GetCustomerFavoriteProducts(c *gin.Context) {
	customerID := c.Param("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Customer ID is required"})
		return
	}

	products, err := h.favoriteService.GetCustomerFavoriteProducts(c.Request.Context(), customerID)
	if err != nil {
		log.Println("Error getting favorite products", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch favorite products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
