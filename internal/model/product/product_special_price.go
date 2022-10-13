package product

import (
	"linebot/internal/config/db"

	"gorm.io/gorm"
)

type Price struct {
	ProductId int
	WeekDay   float32 `gorm:"not null;default:0"`
	Holiday   float32 `gorm:"not null;default:0"`
}

func (price Price) Add() error {
	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Create(&price).Error
	})
}

func GetAllPrice() ([]Price, error) {
	var Prices []Price
	err := db.DB.Find(&Prices).Error

	return Prices, err
}

func GetPriceByProductID(productId int64) (Price, error) {
	var GetPrice Price
	err := db.DB.Where("product_id=?", productId).Find(&GetPrice).Error

	return GetPrice, err
}
