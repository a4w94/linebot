package main

import (
	"fmt"
	"linebot/internal/config"
	_ "linebot/internal/config/db/migrate"
	"linebot/internal/model/product"
	"linebot/internal/model/stock"
	"linebot/internal/route"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	//model.InitDbContext()
	//db.InitDbContext()

	TestData()
	GetData()
	ginroute := route.InitRouter()
	fmt.Printf("Address: http://localhost:%s/ \n", config.HttpPort)
	ginroute.Run(":" + config.HttpPort)

	//first page
}

func TestData() {
	var p1 = product.Product{
		CampRoundName: "A區",
		Size:          "5m*5m",
		Price_Weekday: 1000,
		Price_Holiday: 1200,
		Uint:          "帳",
		TotlaNum:      5,
		ImageUri:      []string{"https://i.imgur.com/XXwY96T.jpeg", "https://i.imgur.com/3dthZKo.jpeg"},
	}

	var p2 = product.Product{
		CampRoundName: "B區-3帳包區",
		Size:          "5m*5m",
		Price_Weekday: 1000,
		Price_Holiday: 1200,
		Uint:          "區",
		TotlaNum:      1,
		ImageUri:      []string{"https://i.imgur.com/XXwY96T.jpeg", "https://i.imgur.com/3dthZKo.jpeg"},
	}
	p1.Add()
	p2.Add()
	//all, _ := product.GetAll()
	//fmt.Println(all)

	t, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

	for i := 0; i < 5; i++ {
		t = t.AddDate(0, 0, 1)
		var r_n = 5

		var tmp = stock.Stock{
			Date:      t,
			ProductId: 1,
			RemainNum: r_n,
		}

		var tmp1 = stock.Stock{
			Date:      t,
			ProductId: 2,
			RemainNum: 1,
		}
		tmp.Add()
		tmp1.Add()
	}

	// o1 := order.Order{
	// 	OrderSN:      order.GenerateOrderSN(0),
	// 	UserID:       "jjwwt",
	// 	ProductId:    1,
	// 	Amount:       3,
	// 	PaymentTotal: 1000,
	// }
	// err := o1.Add()
	// if err != nil {
	// 	fmt.Println("insert o1 failed")
	// }
	// o1.Add()
	// if err != nil {
	// 	fmt.Println("insert o1_2 failed")
	// }
	// var check_insert_success bool
	// var index int
	// for !check_insert_success {
	// 	o2 := order.Order{
	// 		OrderSN:      order.GenerateOrderSN(index),
	// 		UserID:       "jjwwt",
	// 		ProductId:    2,
	// 		Amount:       5,
	// 		PaymentTotal: 5000,
	// 	}
	// 	err := o2.Add()
	// 	if err != nil {
	// 		fmt.Println("insert o2 failed")
	// 		index++
	// 	} else {
	// 		check_insert_success = true
	// 	}
	// }

}

func GetData() {

	// var o = order.Order{
	// 	OrderSN:   "dasfe123",
	// 	Amount:    1,
	// 	ProductId: 2,
	// 	UserID:    "U8d3ff666c698729d2de5d62cf4607964",
	// }
	// o.Add()

	// o_s, _ := order.GetAllOrder()
	// for _, r := range o_s {
	// 	fmt.Println(r)
	// }

	// o_s2, _ := order.GetOrdersByUserID("U8d3ff666c698729d2de5d62cf4607964")
	// for _, r := range o_s2 {
	// 	fmt.Println(r)
	//}

	// stocks, _ := stock.GetAll()
	// for _, s := range stocks {
	// 	fmt.Println(s)
	// }

	// start, _ := time.Parse("2006-01-02", time.Now().AddDate(0, 0, 1).Format("2006-01-02"))

	// end, _ := time.Parse("2006-01-02", time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
	// ordernum := 2
	// stocks, _ := stock.GetStocks_By_ID_and_DateRange(1, start, end)
	// fmt.Println("Get id 1 stocks")

	// for _, r := range stocks {
	// 	fmt.Println(r)
	// }
	// for i := range stocks {
	// 	stocks[i].RemainNum -= ordernum
	// 	stock.UpdateStock(stocks[i])
	// }
	// stocks2, _ := stock.GetStocks_By_ID_and_DateRange(1, start, end)
	// fmt.Println("Get id 1 stocks update")

	// for _, r := range stocks2 {
	// 	fmt.Println(r)
	// }

	// ps, _ := product.GetAll()
	// for _, p := range ps {
	// 	fmt.Println(p, p.ID)
	// }
	// var t line.Search_Time
	// t.Start, _ = time.Parse("2006-01-02", "2022-09-25")
	// t.End, _ = time.Parse("2006-01-02", "2022-09-30")
	// r := line.SearchRemainCamp(t)
	// for _, t := range r {
	// 	fmt.Println(t.Product)
	// 	fmt.Println(t.Stocks)
	// 	fmt.Println(t.RemainMinAmount)
	// 	fmt.Println()
	// }

}
