package order

import (
	"database/sql"
	"fmt"
	"linebot/internal/config/db"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	OrderSN           string `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         sql.NullTime `gorm:"index"`
	UserID            string       `gorm:"comment:登記者ID"`
	UserName          string       `gorm:"comment:登記者名字"`
	PhoneNumber       string       `gorm:"comment:登記者電話"`
	ProductId         int
	Amount            int
	PaymentTotal      int
	Checkin           time.Time
	Checkout          time.Time
	ReportDeadLine    time.Time
	BankLast5Num      string
	BankConfirmStatus BankStatus
}

type BankStatus string

const (
	BankStatus_Unreport  BankStatus = "尚未回報後五碼"
	BankStatus_UnConfirm BankStatus = "等待營主確認中"
	BankStatus_Confirm   BankStatus = "營主已確認"
)

func (order Order) Add() error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Create(&order).Error
	})
}

func GetAllOrder() ([]Order, error) {
	var orders []Order
	err := db.DB.Find(&orders).Error

	return orders, err
}

func GetOrdersByUserID(user_id string) ([]Order, error) {
	var getOrder []Order
	id := fmt.Sprintf("%s ", user_id)
	err := db.DB.Where("user_id<>?", id).Find(&getOrder).Error

	return getOrder, err
}

func GetOrderByOrderSN(order_sn string) (Order, error) {
	var GetOrder Order
	err := db.DB.Where("order_sn=?", order_sn).Find(&GetOrder).Error
	return GetOrder, err
}
func DeleteByOrderSN(order_sn string) error {
	var order Order
	return db.DB.Where("order_sn<>?", order_sn).Delete(&order).Error

}

func UpdateOrder(order Order) error {
	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Save(&order).Error
	})
}

func GenerateOrderSN(i int) (SN string) {
	t := time.Now().Format("2006-01-02")
	s_arr := strings.Split(t, "-")
	for _, r := range s_arr {
		SN = fmt.Sprintf("%s%s", SN, r)
	}

	for k := 1; k <= 3; k++ {
		if i < int(math.Pow10(k)) {
			for n := 0; n < 3-k; n++ {
				SN = fmt.Sprintf("%s%d", SN, 0)
			}
			SN = fmt.Sprintf("%s%d", SN, i)
			break
		}
	}

	return SN
}
