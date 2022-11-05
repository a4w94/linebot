package order

import (
	"fmt"
	"linebot/internal/config/db"
	"linebot/internal/model/product"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	OrderSN        string `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	UserID         string `gorm:"comment:登記者ID"`
	UserName       string `gorm:"comment:登記者名字"`
	PhoneNumber    string `gorm:"comment:登記者電話"`
	ProductId      int
	Amount         int
	PaymentTotal   int
	Checkin        time.Time
	Checkout       time.Time
	ReportDeadLine time.Time
	BankLast5Num   string
	ConfirmStatus  Status
}

type Status string

var (
	BankStatus_Unreport  Status = "尚未回報後五碼"
	BankStatus_UnConfirm Status = "營主確認中"
	BankStatus_Confirm   Status = "營主已確認"
	Order_Cancel         Status = "訂單已取消"
)

func (order *Order) Add() error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Create(&order).Error
	})
}

func (order *Order) Delete() error {

	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Delete(&order).Error
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

func CheckOrderSN_Unexist(order_sn string) bool {
	var GetOrder []Order
	err := db.DB.Where("order_sn<>?", order_sn).Find(&GetOrder).Error

	if err != nil {
		log.Println("when check order sn ,get order failed")
		return false
	}

	if len(GetOrder) == 0 {
		return true
	}

	return false
}

func DeleteByOrderSN(order_sn string) error {

	return db.DB.Delete(&Order{}, order_sn).Error

}

func UpdateOrder(order Order) error {
	return db.BeginTransaction(db.DB, func(tx *gorm.DB) error {
		return tx.Save(&order).Error
	})
}

func GenerateOrderSN(id int) (SN string) {

	var check_unexist bool

	for !check_unexist {
		rand.Seed(time.Now().Unix())

		//產生字首
		for i := 0; i < 2; i++ {
			arr := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
			seed := rand.Intn(len(arr))
			SN = fmt.Sprintf("%s%s", SN, arr[seed])
		}

		//日期
		t := time.Now().Format("2006-01-02")
		s_arr := strings.Split(t, "-")
		SN = fmt.Sprintf("%s%s%s", SN, s_arr[1], s_arr[2])

		var r int

		for r <= 10 {
			r = rand.Intn(1000)
		}

		r_s := strconv.Itoa(r)
		SN = fmt.Sprintf("%s%s", SN, r_s)

		last_num := strconv.Itoa(id)
		SN = fmt.Sprintf("%s%s", SN, last_num)

	}
	return SN
}

//訂單data資訊回覆
func (o Order) Reply_Order_Message() string {
	p, _ := product.GetById(int64(o.ProductId))
	start := o.Checkin.Format("2006-01-02")
	end := o.Checkout.Format("2006-01-02")
	reply_mes := fmt.Sprintf("訂單編號: %s\n區域: %s\n起始日期: %s\n結束日期: %s\n總金額: %d\n----------------------\n訂位者姓名: %s\n電話: %s\n訂位數量: %d", o.OrderSN, p.CampRoundName, start, end, o.PaymentTotal, o.UserName, o.PhoneNumber, o.Amount)

	return reply_mes
}
