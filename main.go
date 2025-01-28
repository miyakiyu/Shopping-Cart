package main

import (
	"log"
	"shop/controller"
	db "shop/database"

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

	// route
	r.GET("/ping", controller.GetAllProducts)
	r.POST("/register", func(c *gin.Context) {
		controller.UserRegister(c, database)
	})
	r.POST("/login", func(c *gin.Context) {
		controller.UserLogin(c, database)
	})

	// Test
	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.Run()
}
