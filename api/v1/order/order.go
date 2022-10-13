package order

import (
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

	orders.Add()

}
