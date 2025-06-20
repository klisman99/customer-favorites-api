package middleware

import (
	"app/internal/infra/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header is required"},
			)
			return
		}

		userID, err := tokenService.Validate(token)
		if err != nil {
			log.Println("Error validating token", err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired token"},
			)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
