package line

import (
	"fmt"
	"linebot/internal/model/order"
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

func (t Search_Time) SearchRemainCamp_ALL() (r_c []RemainCamp, err error) {
	fmt.Println("all time:", t)

	products, err := product.GetAll()
	if err != nil {
		log.Println("Get Products Failed", err)
	}

	//fmt.Println("開始搜尋全部剩餘營位")
	for _, p := range products {
		tmp, err := t.SearchRemainCamp(p)
		if err != nil {
			return r_c, err
		}

		r_c = append(r_c, tmp)
	}
	fmt.Println("rc", r_c)

	return r_c, err
}

func (t Search_Time) SearchRemainCamp(p product.Product) (tmp RemainCamp, err error) {

	tmp.Product = p
	var err1 error
	tmp.Stocks, err1 = stock.GetStocks_By_ID_and_DateRange(tmp.Product.ID, t.Start, t.End)
	if err1 != nil {
		log.Println("GetStocks Failed", err1)
		return tmp, err1
	}

	tmp.PaymentTotal = t.camp_PaymentTotal(p)

	var remain []int
	for _, s := range tmp.Stocks {

		remain = append(remain, s.RemainNum)
	}

	//找到最小剩餘數
	tmp.RemainMinAmount, _ = tool.Find_Min_and_Max(remain)

	//加總總金額

	return tmp, err
}

func (t Search_Time) camp_PaymentTotal(p product.Product) (paymentTotal int) {

	for !t.Start.Equal(t.End) {

		paymentTotal += p.Price_Weekday

		paymentTotal += p.Price_Weekday
		t.Start = t.Start.AddDate(0, 0, 1)

	}
	return paymentTotal
}

func (t Search_Time) Check_Remain_Num_Enough(input_order_num int, region_name string) (bool, error) {

	p, err := product.GetIdByCampRoundName(region_name)
	if err != nil {
		log.Fatal("check remain get product failed")
	}

	r, err := t.SearchRemainCamp(p)
	if err != nil {
		return false, err
	}

	remain := r.RemainMinAmount

	return input_order_num <= remain, err

}

func (t Search_Time) Update_Stock_Remain_by_Order(o order.Order) {
	stocks, _ := stock.GetStocks_By_ID_and_DateRange(uint(o.ProductId), t.Start, t.End)

	for i := range stocks {
		stocks[i].RemainNum -= o.Amount
		stocks[i].UpdateStock()
	}
}
