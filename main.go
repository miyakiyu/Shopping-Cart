package main

import (
	"log"
	"net/http"
	"shop/controller"
	db "shop/database"
	"shop/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// database
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("連線錯誤: ", err)
	}
	// test db connection
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("連線錯誤: ", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("連線錯誤", err)
	} else {
		log.Println("資料庫連線成功")
	}

	//settings
	r := gin.Default()
	r.LoadHTMLGlob("html/*.html")
	r.Use(middleware.AuthMiddleware())

	// route
	r.POST("/register", func(c *gin.Context) {
		controller.UserRegister(c, database)
	})
	r.POST("/login", func(c *gin.Context) {
		controller.UserLogin(c, database)
	})
	r.POST("/add-product", func(c *gin.Context) {
		controller.AddProduct(c, database)
	})
	r.GET("/user-info", func(c *gin.Context) {
		controller.UserInfo(c)
	})
	r.GET("/products", func(c *gin.Context) {
		controller.GetAllProducts(c, database)
	})
	r.POST("/add-to-cart", func(c *gin.Context) {
		controller.AddToCart(c, database)
	})
	r.POST("/remove", func(c *gin.Context) {
		controller.RemoveFromCart(c, database)
	})
	r.GET("/shopcart", func(c *gin.Context) {
		controller.GetCartItem(c, database)
	})
	r.POST("/checkout-cart", func(c *gin.Context) {
		controller.Checkout(c, database)
	})

	// HTML
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/add-product", func(c *gin.Context) {
		// Check user role
		role, exists := c.Get("role")
		if !exists || role != "seller" {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusOK, "add_product.html", nil)
	})
	r.GET("/shop", func(c *gin.Context) {
		c.HTML(http.StatusOK, "shop.html", nil)
	})
	r.GET("/cart", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cart.html", nil)
	})
	r.GET("/checkout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "checkout.html", nil)
	})

	r.Run()
}
