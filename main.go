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

	// route
	r := gin.Default()

	r.GET("/ping", controller.GetAllProducts)
	r.POST("/register", func(c *gin.Context) {
		controller.RegisterUser(c, database)
	})

	r.Run()
}
