package line

import (
	"fmt"
	"linebot/internal/model/product"
	"linebot/internal/model/stock"
	"linebot/pkg/tool"
	"log"
)

type RemainCamp struct {
	Product         product.Product
	Stocks          []stock.Stock
	RemainMinAmount int
	PaymentTotal    int
}

func (t Search_Time) SearchRemainCamp_ALL() (r_c []RemainCamp) {
	fmt.Println("all time:", t)
	var err error
	products, err := product.GetAll()
	if err != nil {
		log.Println("Get Products Failed", err)
	}

	fmt.Println("開始搜尋全部剩餘營位")
	for _, p := range products {
		tmp := t.SearchRemainCamp(p)

		r_c = append(r_c, tmp)
	}

	return r_c
}

func (t Search_Time) SearchRemainCamp(p product.Product) (tmp RemainCamp) {
	fmt.Println("search camp input time:", t)
	fmt.Println("input product:", p)
	tmp.Product = p
	tmp.Stocks, err = stock.GetStocks_By_ID_and_DateRange(tmp.Product.ID, t.Start, t.End)
	if err != nil {
		log.Println("GetStocks Failed", err)
	}
	fmt.Println("stocks", tmp.Stocks)

	tmp.PaymentTotal = t.camp_PaymentTotal(p)
	fmt.Println("pay:", tmp.PaymentTotal)
	var remain []int
	for _, s := range tmp.Stocks {
		fmt.Println("stock", s)
		remain = append(remain, s.RemainNum)
	}
	fmt.Println("remain arr", remain)
	//找到最小剩餘數
	tmp.RemainMinAmount, _ = tool.Find_Min_and_Max(remain)

	//加總總金額

	return tmp
}

func (t Search_Time) camp_PaymentTotal(p product.Product) (paymentTotal int) {

	for {
		if t.Start.Equal(t.End) {
			paymentTotal += p.Price_Weekday
			break
		}

		paymentTotal += p.Price_Weekday
		t.Start = t.Start.AddDate(0, 0, 1)

	}
	return paymentTotal
}

func (t Search_Time) Check_Remain_Num_Enough(input_order_num int, region_name string) bool {
	fmt.Println("check remain", input_order_num, region_name)
	fmt.Println("input time", t)
	p, err := product.GetIdByCampRoundName(region_name)
	if err != nil {
		log.Fatal("check remain get product failed")
	}
	fmt.Println("product", p)
	remain := t.SearchRemainCamp(p).RemainMinAmount
	fmt.Println("remain", remain)
	return input_order_num <= remain

}
