package middlewares

import (
	"net/http"

	"ashwin.com/go-auth/helper"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Your not authorized")
			c.Abort()
			return
		}
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, msg := helper.ValidateToken(clientToken)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}

		c.Set("name", claims.Name)
		c.Next()
	}
}
