package main

import (
	"fmt"
	"linebot/api/v1/line"
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
		ImageUri:      []string{"https://i.imgur.com/XXwY96T.jpeg", "https://i.imgur.com/3dthZKo.jpeg"},
	}

	var p2 = product.Product{
		CampRoundName: "B區",
		Size:          "5m*5m",
		ImageUri:      []string{"https://i.imgur.com/XXwY96T.jpeg", "https://i.imgur.com/3dthZKo.jpeg"},
	}
	p1.Add()
	p2.Add()
	//all, _ := product.GetAll()
	//fmt.Println(all)

	t, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

	for i := 0; i < 20; i++ {
		t = t.AddDate(0, 0, 1)
		var r_n = 5

		var tmp = stock.Stock{
			Date:      t,
			ProductId: 1,
			TotlaNum:  5,
			RemainNum: r_n,
		}

		var tmp1 = stock.Stock{
			Date:      t,
			ProductId: 2,
			TotlaNum:  5,
			RemainNum: r_n,
		}
		tmp.Add()
		tmp1.Add()
	}

}

func GetData() {

	stocks, _ := stock.GetAll()
	for _, s := range stocks {
		fmt.Println(s)
	}
	ps, _ := product.GetAll()
	for _, p := range ps {
		fmt.Println(p, p.ID)
	}
	var t line.Search_Time
	t.Start, _ = time.Parse("2006-01-02", "2022-09-25")
	t.End, _ = time.Parse("2006-01-02", "2022-09-30")
	r := line.SearchRemainCamp(t)
	for _, t := range r {
		fmt.Println(t.Product)
		fmt.Println(t.Stocks)
		fmt.Println(t.RemainMinAmount)
		fmt.Println()
	}
}
