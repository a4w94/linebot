package order

import (
	"fmt"
	"linebot/internal/model/order"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var orders = order.Order{
		OrderSN:     "GG123",
		UserID:      "k123",
		UserName:    "JJ",
		PhoneNumber: "0909990",
	}
	var item = order.OrderItem{
		OrderSN:      orders.OrderSN,
		ProductId:    1,
		Amount:       9,
		PaymentTotal: 1000,
	}
	orders.Add()
	item.Add()
	allorder, _ := order.GetAllOrder()
	for _, r := range allorder {
		fmt.Println(r)
		it, _ := order.GetOrderItemByOrderSN(r.OrderSN)
		fmt.Println(it)
	}

}
