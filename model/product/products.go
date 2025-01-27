package products

type Product struct {
	ID    int    `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Price int    `gorm:"not null"`
	Stock int    `gorm:"not null"`
}
