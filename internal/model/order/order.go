package order

import (
	"linebot/internal/config/db"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderSN     string `gorm:"comment:訂單序號"`
	UserID      string `gorm:"comment:登記者ID"`
	UserName    string `gorm:"comment:登記者名字"`
	PhoneNumber string `gorm:"comment:登記者電話"`
}

func (order Order) Add() error {
	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Create(&order).Error
	})
}

func GetAllOrder() ([]Order, error) {
	var orders []Order
	err := db.DB.Find(&orders).Error

	return orders, err
}

func GetOrderByUserID(user_id int64) (Order, error) {
	var getOrder Order
	err := db.DB.Where("user_id=?", user_id).Find(&getOrder).Error

	return getOrder, err
}
