package main

import (
	_ "app/docs"
	"app/internal/api/handler"
	"app/internal/api/middleware"
	domainservice "app/internal/domain/service"
	"app/internal/infra/db"
	infraservice "app/internal/infra/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Customer Favorites API
// @version 1.0
// @description This is an API to manage customer favorites products

// @host localhost:3002
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	database := db.Connect()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_jwt_secret"
	}

	userRepository := db.NewUserRepository(database)
	tokenService := infraservice.NewTokenService(jwtSecret, time.Hour)
	authService := domainservice.NewAuthService(userRepository, tokenService)
	authHandler := handler.NewAuthHandler(authService)

	customerRepository := db.NewCustomerRepository(database)
	customerService := domainservice.NewCustomerService(customerRepository)
	customerHandler := handler.NewCustomerHandler(customerService)

	favoriteRepository := db.NewFavoriteRepository(database)
	productService := infraservice.NewProductService()
	favoriteService := domainservice.NewFavoriteService(favoriteRepository, productService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	router := gin.Default()

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		if err := database.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	router.POST("/signup", authHandler.SignUp)
	router.POST("/signin", authHandler.SignIn)

	v1Api := router.Group("/api/v1")
	v1Api.Use(middleware.AuthMiddleware(tokenService))
	{
		customers := v1Api.Group("/customers")
		{
			customers.POST("", customerHandler.Create)
			customers.GET("/:customer_id", customerHandler.GetByID)
			customers.GET("", customerHandler.GetAll)
			customers.PUT("/:customer_id", customerHandler.Update)
			customers.DELETE("/:customer_id", customerHandler.Delete)

			favorites := customers.Group("/:customer_id/favorites")
			{
				favorites.POST("", favoriteHandler.AddFavorite)
				favorites.GET("", favoriteHandler.GetCustomerFavoriteProducts)
				favorites.DELETE("/:product_id", favoriteHandler.RemoveFavorite)
			}
		}
	}

	if err := router.Run(":3002"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
