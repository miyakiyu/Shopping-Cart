package db

import (
	products "shop/model/product"
	"shop/model/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=miyakiyu password=12345 dbname=miyakiyu port=5432 sslmode=disable TimeZone=Asia/Taipei"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&products.Product{})
	db.AutoMigrate(&user.User{})
	return db, nil
}
