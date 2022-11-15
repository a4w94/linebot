package main

import (
	"fmt"
	"linebot/internal/config"
	_ "linebot/internal/config/db/migrate"
	"linebot/internal/model/product"
	"linebot/internal/model/stock"
	"linebot/internal/richmenu"
	"linebot/internal/route"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	//model.InitDbContext()
	//db.InitDbContext()
	richmenu.Build_RichMenu()
	//TestData()
	//GetData()
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

	for i := 0; i < 4; i++ {
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

	// start, _ := time.Parse("2006-01-02", "2022-11-08")
	// end, _ := time.Parse("2006-01-02", "2022-11-08")
	// var t line.Search_Time
	// t.Start = start
	// t.End = end
	// fmt.Println(stock.GetStocks_By_ID_and_DateRange(1, start, end))

	//fmt.Println(order.CheckOrderSN_exist("RO110779510"))
	// s, err := stock.GetStocks_By_ID_and_DateRange(1, start, end)
	// for _, r := range s {
	// 	fmt.Println(r)
	// }
	// fmt.Println(err)

	// o_s, _ := order.GetAllOrder()
	// for _, r := range o_s {
	// 	fmt.Println(r.UserID, r)
	// }

	// line.Carousel_Orders(o_s2)
	// o, err := order.GetOrderByOrderSN("20221103000")
	// fmt.Println("get", o, err)

	// // stocks, _ := stock.GetAll()
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
