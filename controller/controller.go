package controller

import (
	"net/http"
	"shop/model/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User register
func RegisterUser(c *gin.Context, db *gorm.DB) {
	var input struct {
		Username string `json:"user" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" bindng:"required,oneof = buyer seller"`
	}
	// Check input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Store user to database
	user := user.User{
		Username: input.Username,
		Password: input.Password,
		Role:     input.Role,
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created!"})
}

// Show all products
func GetAllProducts(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
