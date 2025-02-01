package cart

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
