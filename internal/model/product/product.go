package product

import (
	"linebot/internal/config/db"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CampRoundName string `gorm:"not null"`
	Price         Price  `gorm:"embedded"`
	Size          string
	Image         int
	Description   string
}

func AddN(input interface{}) error {
	return db.BeginTranscation(db.DB, func(tx *gorm.DB) error {
		return tx.Create(&input).Error
	})
}

func Add(product *Product) error {
	return db.BeginTranscation(db.DB, func(tx *gorm.DB) error {
		return tx.Create(&product).Error
	})
}

func GetAll() ([]Product, error) {
	var Products []Product
	err := db.DB.Find(&Products).Error

	return Products, err
}

func GetById(Id int64) (Product, error) {
	var GetProduct Product
	err := db.DB.Where("Id=?", Id).Find(&GetProduct).Error

	return GetProduct, err
}

func DeleteById(Id int64) (Product, error) {
	var product Product
	err := db.DB.Where("Id=?", Id).Delete(&product).Error
	return product, err
}
