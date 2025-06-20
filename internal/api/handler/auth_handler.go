package handler

import (
	"app/internal/domain"
	"app/internal/domain/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service *service.AuthService
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// @Summary Sign up
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body AuthRequest true "Sign up credentials"
// @Success 201 "Created"
// @Failure 400 "Invalid request data"
// @Failure 409 "Username already exists"
// @Failure 500 "Failed to sign up"
// @Router /signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req AuthRequest
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

	_, err := h.service.SignUp(c, req.Username, req.Password)
	if err != nil {
		if err == domain.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"message": "Username already exists"})
			return
		}
		log.Println("Error signing up", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to sign up"})
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Sign in
// @Description Authenticate a user and return a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body AuthRequest true "Sign in credentials"
// @Success 200 "OK"
// @Failure 400 "Invalid request data"
// @Failure 401 "Unauthorized"
// @Failure 500 "Failed to sign in"
// @Router /signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req AuthRequest
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

	token, err := h.service.SignIn(c, req.Username, req.Password)

	if err != nil {
		log.Println("Error signing in", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to sign in"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
