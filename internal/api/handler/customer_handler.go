package handler

import (
	"app/internal/domain"
	"app/internal/domain/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	service *service.CustomerService
}

type CustomerCreateRequest struct {
	Name  string `json:"name" validate:"required,min=3" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

type CustomerUpdateRequest struct {
	Name  string `json:"name" validate:"required,min=3" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

type CustomerResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Error message"`
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// @Summary Create a new customer
// @Description Creates a new customer with the provided details
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body CustomerCreateRequest true "Customer details"
// @Success 201 "Created customer details"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 409 {object} ErrorResponse "Email already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var customerCreateRequest CustomerCreateRequest
	if err := c.ShouldBindJSON(&customerCreateRequest); err != nil {
		log.Println("Error binding request data", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	if err := validator.New().Struct(customerCreateRequest); err != nil {
		log.Println("Error validating request data", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	_, err := h.service.Create(c.Request.Context(), customerCreateRequest.Name, customerCreateRequest.Email)
	if err != nil {
		if err == domain.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"message": "Email already exists"})
			return
		}
		log.Println("Error creating customer", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create customer"})
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Get customer by ID
// @Description Retrieves a customer by their unique identifier
// @Tags Customer
// @Produce json
// @Param customer_id path string true "Customer ID" example="550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} CustomerResponse "Customer details"
// @Failure 400 {object} ErrorResponse "Invalid customer ID format"
// @Failure 404 {object} ErrorResponse "Customer not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/customers/{customer_id} [get]
func (h *CustomerHandler) GetByID(c *gin.Context) {
	customerID := c.Param("customer_id")
	if _, err := uuid.Parse(customerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid customer ID"})
		return
	}

	customer, err := h.service.GetByID(c.Request.Context(), customerID)
	if err != nil {
		log.Println("Error fetching customer", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch customer"})
		return
	}

	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// @Summary Get all customers
// @Description Retrieves a list of all customers
// @Tags Customer
// @Produce json
// @Success 200 {array} CustomerResponse "List of customers"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/customers [get]
func (h *CustomerHandler) GetAll(c *gin.Context) {
	customers, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		log.Println("Error fetching customers", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch customers"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// @Summary Update customer
// @Description Updates an existing customer's details
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer_id path string true "Customer ID" example="550e8400-e29b-41d4-a716-446655440000"
// @Param customer body CustomerUpdateRequest true "Updated customer details"
// @Success 204 "No content"
// @Failure 400 {object} ErrorResponse "Invalid request data or customer ID"
// @Failure 404 {object} ErrorResponse "Customer not found"
// @Failure 409 {object} ErrorResponse "Email already in use"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/customers/{customer_id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	customerID := c.Param("customer_id")
	if _, err := uuid.Parse(customerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid customer ID"})
		return
	}

	var customerUpdateRequest CustomerUpdateRequest
	if err := c.ShouldBindJSON(&customerUpdateRequest); err != nil {
		log.Println("Error binding request data", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	if err := validator.New().Struct(customerUpdateRequest); err != nil {
		log.Println("Error validating request data", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	err := h.service.Update(c.Request.Context(), customerID, customerUpdateRequest.Name, customerUpdateRequest.Email)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
			return
		}
		if err == domain.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"message": "Email already exists"})
			return
		}
		log.Println("Error updating customer", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update customer"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Delete customer
// @Description Deletes a customer by their ID
// @Tags Customer
// @Param customer_id path string true "Customer ID" example="550e8400-e29b-41d4-a716-446655440000"
// @Success 204 "No content"
// @Failure 400 {object} ErrorResponse "Invalid customer ID format"
// @Failure 404 {object} ErrorResponse "Customer not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/customers/{customer_id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	customerID := c.Param("customer_id")
	if _, err := uuid.Parse(customerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid customer ID"})
		return
	}

	err := h.service.Delete(c.Request.Context(), customerID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
			return
		}
		log.Println("Error deleting customer", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete customer"})
		return
	}

	c.Status(http.StatusNoContent)
}
