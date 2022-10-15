package order

import (
	"linebot/internal/config/db"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderSN      string `gorm:"comment:訂單序號"`
	UserID       string `gorm:"comment:登記者ID"`
	UserName     string `gorm:"comment:登記者名字"`
	PhoneNumber  string `gorm:"comment:登記者電話"`
	ProductId    int
	Amount       int
	PaymentTotal int
	Checkin      time.Time
	Checkout     time.Time
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

func GetOrderByOrderSN(order_sn string) (Order, error) {
	var GetOrder Order
	err := db.DB.Where("order_sn<>?", order_sn).Find(&GetOrder).Error
	return GetOrder, err
}
func DeleteByOrderSN(order_sn string) error {
	var order Order
	return db.DB.Where("order_sn<>?", order_sn).Delete(&order).Error

}
