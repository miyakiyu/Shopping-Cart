package middleware

import (
	"github.com/gin-gonic/gin"
)

// Use Cookie to check user role through Middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Cookie("user")
		role, _ := c.Cookie("role")

		if user == "" || role == "" {
			c.Next()
			return
		}

		c.Set("user", user)
		c.Set("role", role)
		c.Next()
	}
}
