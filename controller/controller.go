package controller

import (
	"net/http"
	products "shop/model/product"
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
	//Setting cookie
	c.SetCookie("user", user.Username, 3600, "/", "localhost", false, true)
	c.SetCookie("role", user.Role, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login success!"})
}

// Add new product
func AddProduct(c *gin.Context, db *gorm.DB) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Price int    `json:"price" binding:"required"`
		Stock int    `json:"stock" binding:"required"`
	}
	// Check input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不是json"})
		return
	}
	// Check role
	role, exist := c.Get("role")
	if !exist || role != "seller" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}
	// Product to database
	product := products.Product{
		Name:  input.Name,
		Price: input.Price,
		Stock: input.Stock,
	}
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database錯誤"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product added!"})

}

// Check user role
func UserInfo(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"logged_in": true,
		"username":  user.(string),
		"role":      c.GetString("role"),
	})
}

// Show all products from database
func GetAllProducts(c *gin.Context, db *gorm.DB) {

}
