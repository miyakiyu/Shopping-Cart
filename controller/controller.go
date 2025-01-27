package controller

import (
	"net/http"
	"shop/model/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// User register
func UserRegister(c *gin.Context, db *gorm.DB) {
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

// User Login
func UserLogin(c *gin.Context, db *gorm.DB) {
	var input struct {
		Username string `json:"user" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	// Check input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Find user in database
	var user user.User
	if err := db.Where("username = ? AND password = ?", input.Username, input.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login success!"})
}

// Show all products
func GetAllProducts(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
