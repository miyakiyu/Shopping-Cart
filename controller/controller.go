package controller

import (
	"net/http"
	"shop/model/cart"
	products "shop/model/product"
	"shop/model/user"
	"strconv"

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
	c.SetCookie("user", user.Username, 3600, "/", "localhost", false, false)
	c.SetCookie("role", user.Role, 3600, "/", "localhost", false, false)
	c.SetCookie("user_id", strconv.Itoa(user.ID), 3600, "/", "localhost", false, false)

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
	var list []products.Product
	if err := db.Find(&list).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func GetCartItem(c *gin.Context, db *gorm.DB) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	// Trans string(cookie) to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User_id is not int"})
		return
	}
	var list []cart.Cart
	if err := db.Preload("Product").Where("user_id = ?", userIDInt).Find(&list).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Add product to shopping cart
func AddToCart(c *gin.Context, db *gorm.DB) {
	var request struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	var item cart.Cart
	err := db.Where("user_id = ? AND product_id = ?", request.UserID, request.ProductID).First(&item).Error

	if err == nil {
		item.Quantity = request.Quantity + item.Quantity
		if err := db.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
			return
		}
	} else if err == gorm.ErrRecordNotFound {
		newItem := cart.Cart{
			UserID:    request.UserID,
			ProductID: request.ProductID,
			Quantity:  request.Quantity,
		}
		if err := db.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart!"})
}

// Remove product from shopping cart
func RemoveFromCart(c *gin.Context, db *gorm.DB) {
	var request struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Find item in cart
	var item cart.Cart
	if err := db.Where("user_id = ? AND product_id = ?", request.UserID, request.ProductID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// If Quantity less than 0 , delete item
	if item.Quantity-request.Quantity <= 0 {
		db.Unscoped().Where("user_id = ? AND product_id = ?", request.UserID, request.ProductID).Delete(&cart.Cart{})
		c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
		return
	}

	// else reduce quantity
	if err := db.Model(&item).Where("user_id = ? AND product_id = ?", request.UserID, request.ProductID).
		Update("quantity", gorm.Expr("quantity - ?", request.Quantity)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quantity"})
		return
	}
}

// Checkout shopping cart
func Checkout(c *gin.Context, db *gorm.DB) {
	type request struct {
		UserID int `json:"user_id"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Checking item
	var list []cart.Cart
	if err := db.Preload("Product").Where("user_id = ?", req.UserID).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart items"})
		return
	}

	// Delete stock
	// Start transaction
	deal := db.Begin()
	for _, item := range list {
		// If stock isn't enough, rollback
		if item.Quantity > item.Product.Stock {
			deal.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock not enough"})
			return
		}
		// If stock is enough, update stock
		if err := deal.Model(&products.Product{}).
			Where("id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			deal.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
			return
		}
	}

	// Delete cart item
	if err := deal.Where("user_id = ?", req.UserID).Delete(&cart.Cart{}).Error; err != nil {
		deal.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clan cart"})
		return
	}

	deal.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Checkout success!"})
}
