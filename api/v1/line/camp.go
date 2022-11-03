package line

import (
	"fmt"
	"linebot/internal/model/order"
	"linebot/internal/model/product"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Search_Time struct {
	Start time.Time
	End   time.Time
}
type ParseData struct {
	Action string
	Type   string
	Status string
	Data   string
}

type Order_Info struct {
	Region       string `tag:"區域"`
	Start        string `tag:"起始日期"`
	End          string `tag:"結束日期"`
	UserName     string `tag:"訂位者姓名"`
	PhoneNumber  string `tag:"電話"`
	Amount       string `tag:"訂位數量"`
	PaymentTotal string `tag:"總金額"`
}

var Search map[string]*Search_Time

func init() {
	Search = make(map[string]*Search_Time)

	fmt.Println("Init Search ", Search)

	//richmenu.Build_RichMenu()

}

func CampReply(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(500)
		}
		return
	}

	for _, event := range events {

		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				text_trimspace := strings.TrimSpace(message.Text)

				switch {
				case text_trimspace == "我要訂位":
					reply_date_limit(bot, event)

				case text_trimspace == "營地位置":
					bot.ReplyMessage(event.ReplyToken, linebot.NewLocationMessage("小路露營區", "426台中市新社區崑山里食水嵙6-2號", 24.2402679, 120.7943069)).Do()

				case text_trimspace == "營地資訊":
					Img_Carousel_CampRound_Info(bot, event)

				case text_trimspace == "我的訂單":
					reply_User_All_Orders(bot, event)

				case strings.Contains(text_trimspace, "訂單資訊"):
					get_User_Place(text_trimspace, bot, event)
				case strings.Contains(text_trimspace, "*回報資訊*"):
					fmt.Println(text_trimspace)
					report_Bank_Last_FiveNumbers(bot, event, text_trimspace)

				}
			}
		}

		if event.Type == linebot.EventTypePostback {
			data := Parase_postback(event.Postback.Data)
			switch data.Action {
			case "search":

				switch data.Type {

				case "get_start_time":
					date := event.Postback.Params.Date
					// str := fmt.Sprintf("起始日期:%s", date)
					fmt.Println("get start time", date)
					Search[event.Source.UserID].Start, _ = time.Parse("2006-01-02", date)

					reply_date_limit(bot, event)

				case "get_end_time":
					date := event.Postback.Params.Date

					Search[event.Source.UserID].End, _ = time.Parse("2006-01-02", date)
					fmt.Println("Start Time", Search[event.Source.UserID].Start)
					fmt.Println("End Time", Search[event.Source.UserID].End)

					reply_date_limit(bot, event)

				case "start_search":
					if !Search[event.Source.UserID].Start.IsZero() && !Search[event.Source.UserID].End.IsZero() {
						Camp_Search_Remain(bot, event, *Search[event.Source.UserID])

					}
				}

			case "order":
				switch data.Type {
				case "place":
					data.reply_Order_Confirm(bot, event)

				}

			}
			//fmt.Println("data", event.Postback.Data)

		}
	}

}

//快速回覆營位分區名稱
func Quick_Reply_CampRoundName(bot *linebot.Client, event *linebot.Event) {
	var q_p []*linebot.QuickReplyButton
	products, _ := product.GetAll()
	for _, p := range products {
		tmp := &linebot.QuickReplyButton{
			Action: &linebot.MessageAction{
				Label: p.CampRoundName,
				Text:  p.CampRoundName,
			},
		}
		q_p = append(q_p, tmp)
	}
	fmt.Println("Quick_Reply_CampRoundName", q_p)
	bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("選擇分區").WithQuickReplies(&linebot.QuickReplyItems{
		Items: q_p,
	})).Do()

}

//確認輸入營位分區名是否存在
func Is_Name_Exist(name string) (product.Product, bool) {
	fmt.Println("Is_Name_Exist 輸入", name)
	products, _ := product.GetAll()
	var tmp product.Product
	for _, p := range products {
		if p.CampRoundName == name {
			tmp = p
			fmt.Println("名稱存在")
			return tmp, true
		}
	}
	return tmp, false
}

//搜尋剩餘營位
func Camp_Search_Remain(bot *linebot.Client, event *linebot.Event, t Search_Time) {
	var c_t []*linebot.CarouselColumn
	camp_searchs := t.SearchRemainCamp_ALL()
	// fmt.Println("input search time ", t)
	// for i, r := range camp_searchs {
	// 	fmt.Println(i, ":", r.Stocks)
	// }
	for _, s := range camp_searchs {
		total_num := fmt.Sprintf("總共 %d %s", s.Product.TotlaNum, s.Product.Uint)
		remain_num := fmt.Sprintf("剩餘 %d %s", s.RemainMinAmount, s.Product.Uint)
		start := t.Start.Format("2006-01-02")
		end := t.End.Format("2006-01-02")
		des := fmt.Sprintf("%s ~ %s\n\n每%s NT$%d\n%s\n%s\n", start, end, s.Product.Uint, s.PaymentTotal, total_num, remain_num)
		tmp := linebot.CarouselColumn{
			ThumbnailImageURL:    s.Product.ImageUri[0],
			ImageBackgroundColor: "#000000",
			Title:                s.Product.CampRoundName,
			Text:                 des,
			Actions: []linebot.TemplateAction{
				&linebot.PostbackAction{
					Label:       "我要訂位",
					Data:        fmt.Sprintf("action=order&item=%d&num=%d", s.Product.ID, s.RemainMinAmount),
					InputOption: linebot.InputOptionOpenKeyboard,
					FillInText:  fmt.Sprintf("訂單資訊 \n----------------------\n區域: %s\n起始日期: %s\n結束日期: %s\n----------------------\n訂位者姓名: \n電話: 09\n訂位數量: ", s.Product.CampRoundName, start, end),
				},
			},
		}
		fmt.Println(tmp)
		c_t = append(c_t, &tmp)
	}
	delete(Search, event.Source.UserID)

	bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Camp Search",
		&linebot.CarouselTemplate{
			Columns: c_t,

			ImageAspectRatio: "rectangle",
			ImageSize:        "cover",
		})).Do()

}

func Add_Carousel_Imgae() (c_i []*linebot.ImageCarouselColumn) {
	c1 := linebot.ImageCarouselColumn{
		ImageURL: "https://example.com/bot/images/item1.jpg",
		Action: &linebot.PostbackAction{
			Label: "A區",
			Data:  "action=click&itemid=0",
		},
	}

	c2 := linebot.ImageCarouselColumn{
		ImageURL: "https://example.com/bot/images/item1.jpg",
		Action: &linebot.PostbackAction{
			Label: "B區",
			Data:  "action=click&itemid=1",
		},
	}
	c_i = append(c_i, &c1, &c2)
	return c_i
}

func Img_Carousel_CampRound_Info(bot *linebot.Client, event *linebot.Event) {
	fmt.Println("Img_Carousel_CampRound_Info")
	products, err := product.GetAll()
	if err != nil {
		log.Println("Img_Carousel_CampRound_Info Get All Products Failed")
	}
	var c_t []*linebot.ImageCarouselColumn
	for _, r := range products {

		c1 := linebot.ImageCarouselColumn{
			ImageURL: r.ImageUri[0],
			Action: &linebot.PostbackAction{
				Label: r.CampRoundName,
				Data:  "action=click&itemid=0",
			},
		}

		c_t = append(c_t, &c1)

	}
	bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("img carousel",
		&linebot.ImageCarouselTemplate{
			Columns: c_t,
		})).Do()

}

func reply_date_limit(bot *linebot.Client, event *linebot.Event) {

	value, isExist := Search[event.Source.UserID]
	if !isExist {
		Search[event.Source.UserID] = &Search_Time{}
	}
	var (
		start_time string
		start_init string
		start_min  string
		start_Max  string

		end_time string
		end_init string
		end_min  string
		end_Max  string
	)
	init := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	Max := time.Now().AddDate(0, 0, 365).Format("2006-01-02")
	start_init = init
	start_min = init
	end_init = init
	end_min = init
	start_Max = Max
	end_Max = Max
	start_time = "起始日期 "
	end_time = "結束日期 "

	if value != nil {

		switch {
		case !value.Start.IsZero() && value.End.IsZero():
			date := value.Start.Format("2006-01-02")
			start_time = fmt.Sprintf("起始日期 %s", date)
			end_time = "結束日期 "
			end_init = date
			end_min = date
		case value.Start.IsZero() && !value.End.IsZero():
			date := value.End.Format("2006-01-02")
			start_time = "起始日期 "
			end_time = fmt.Sprintf("結束日期 %s", date)
			start_Max = date
		case !value.Start.IsZero() && !value.End.IsZero():
			date_start := value.Start.Format("2006-01-02")
			date_end := value.End.Format("2006-01-02")

			start_time = fmt.Sprintf("起始日期 %s", date_start)
			end_time = fmt.Sprintf("結束日期 %s", date_end)

			end_init = date_start
			end_min = date_start
			start_Max = date_end
		}
	}
	bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("訂位日期", &linebot.ButtonsTemplate{

		Text: "選擇訂位日期",
		Actions: []linebot.TemplateAction{
			&linebot.DatetimePickerAction{
				Label:   start_time,
				Data:    "action=search&type=get_start_time",
				Mode:    "date",
				Initial: start_init,
				Min:     start_min,
				Max:     start_Max,
			},
			&linebot.DatetimePickerAction{
				Label:   end_time,
				Data:    "action=search&type=get_end_time",
				Mode:    "date",
				Initial: end_init,
				Min:     end_min,
				Max:     end_Max,
			},
			&linebot.PostbackAction{
				Label: "查詢",
				Data:  "action=search&type=start_search",
			},
		},
	})).Do()

}

func parase_Order_Info(info string) (bool, string, Order_Info) {
	// fmt.Println("收到訂單資訊", info)
	var tmp Order_Info
	split := strings.Split(info, "\n")
	info_map := make(map[string]string)

	for _, r := range split {
		if strings.Contains(r, ":") {
			arr := strings.Split(r, ":")
			if strings.TrimSpace(arr[1]) != "" {
				info_map[strings.TrimSpace(arr[0])] = strings.TrimSpace(arr[1])
			}
			// if arr[1] != "" {
			// 	info_map[arr[0]] = arr[1]
			// }
		}
	}
	var s_t Search_Time

	t := reflect.TypeOf(tmp)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("tag")
		fmt.Println(tag)
		if v, ok := info_map[tag]; ok {
			if strings.TrimSpace(v) == "" {
				return false, fmt.Sprintf("%s輸入有誤,請重新訂位", tag), Order_Info{}
			} else {
				switch tag {
				case "區域":
					tmp.Region = v
				case "起始日期":
					s_t.Start, err = time.Parse("2006-01-02", v)
					if err != nil {
						return false, fmt.Sprintf("%s輸入有誤,請重新訂位", tag), Order_Info{}
					}
					tmp.Start = v

				case "結束日期":
					s_t.End, err = time.Parse("2006-01-02", v)
					if err != nil {
						return false, fmt.Sprintf("%s輸入有誤,請重新訂位", tag), Order_Info{}
					}
					tmp.End = v

				case "訂位者姓名":
					tmp.UserName = v
				case "電話":
					if len(v) != 10 {
						return false, "電話輸入有誤，請重新訂位", Order_Info{}
					}
					tmp.PhoneNumber = v
				case "訂位數量":
					num, _ := strconv.Atoi(strings.TrimSpace(v))
					if num == 0 {
						return false, "訂位數量不得為零,請重新訂位", Order_Info{}
					}
					//fmt.Println("s_t", s_t)
					if !s_t.Check_Remain_Num_Enough(num, tmp.Region) {
						return false, "剩餘數量不足,請重新訂位", Order_Info{}
					} else {
						tmp.Amount = v
					}

				}
			}
		}
	}

	//計算訂單金額
	amount, _ := strconv.Atoi(tmp.Amount)
	p, _ := product.GetIdByCampRoundName(tmp.Region)
	pay := s_t.camp_PaymentTotal(p)
	pay *= amount
	tmp.PaymentTotal = strconv.Itoa(pay)

	order_info := fmt.Sprintf("確認訂位資訊 \n----------------------\n區域: %s\n起始日期: %s\n結束日期: %s\n總金額: %s\n-----------\n訂位者姓名: %s\n電話: %s\n訂位數量: %s", tmp.Region, tmp.Start, tmp.End, tmp.PaymentTotal, tmp.UserName, tmp.PhoneNumber, tmp.Amount)
	//fmt.Println("order info", order_info)
	return true, order_info, tmp
}

func Parase_postback(data string) (p_d ParseData) {
	str := strings.Split(data, "&")

	for _, p := range str {
		switch {
		case strings.Contains(p, "action"):
			p_d.Action = get_string_data(p)
		case strings.Contains(p, "type"):
			p_d.Type = get_string_data(p)
		case strings.Contains(p, "status"):
			p_d.Status = get_string_data(p)
		case strings.Contains(p, "data"):
			p_d.Data = get_string_data(p)
		}
	}
	return p_d
}

//獲取使用者訂單 回覆確認訊息
func get_User_Place(text string, bot *linebot.Client, event *linebot.Event) {
	ok, check, _ := parase_Order_Info(text)
	//fmt.Println(ok, check)
	data_yes := fmt.Sprintf("action=order&type=place&status=yes&data=%s", check)
	//fmt.Println("data_yes", data_yes)
	if ok {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("確認訂位資訊",
			&linebot.ConfirmTemplate{
				Text: check,
				Actions: []linebot.TemplateAction{
					&linebot.PostbackAction{
						Label: "是",
						Data:  data_yes,
						Text:  "是",
					},
					&linebot.PostbackAction{
						Label: "否",
						Data:  "action=order&type=place&status=no",
						Text:  "否",
					},
				},
			})).Do()
	} else {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(check)).Do()
	}
}

func (p_d ParseData) reply_Order_Confirm(bot *linebot.Client, event *linebot.Event) {
	var reply_mes string

	if p_d.Status == "no" {
		reply_mes = "如有需要請重新訂位，謝謝"
	} else if p_d.Status == "yes" {
		_, _, info := parase_Order_Info(p_d.Data)

		amount, _ := strconv.Atoi(info.Amount)
		product, err := product.GetIdByCampRoundName(info.Region)
		paymenttotal, _ := strconv.Atoi(info.PaymentTotal)
		var order_sn string
		if err != nil {
			fmt.Println("get id failed")
		}
		search_time := parse_string_to_SearchTime(info.Start, info.End)

		var check_insert_success bool
		var index int
		deadline := time.Now().AddDate(0, 0, 3)

		for !check_insert_success {
			order_sn = order.GenerateOrderSN(index)
			var tmp_order = order.Order{
				OrderSN:           order_sn,
				UserID:            event.Source.UserID,
				UserName:          info.UserName,
				PhoneNumber:       info.PhoneNumber,
				ProductId:         int(product.ID),
				Amount:            amount,
				PaymentTotal:      paymenttotal,
				Checkin:           search_time.Start,
				Checkout:          search_time.End,
				ReportDeadLine:    deadline,
				BankConfirmStatus: order.BankStatus_Unreport,
			}

			err = tmp_order.Add()

			// order, _ := order.GetAllOrder()
			// fmt.Println("Order", order)
			if err != nil {
				log.Println("新增訂單失敗", err)
				index++
				reply_mes = "訂位失敗，請重新查詢"
			} else {
				check_insert_success = true
				search_time.Update_Stock_Remain_by_Order(tmp_order)

			}
		}
		remit := fmt.Sprintf("請於%s 23:59前完成匯款並於 *我的訂單* 回報帳號後5碼\n銀行代號: 822\n銀行名稱: 中國信託商業銀行\n匯款帳號: 0342523515\n匯款金額: %s\n", deadline.Format("2006-01-02"), info.PaymentTotal)
		reply_mes = fmt.Sprintf("以下是您的訂位資訊 \n----------------------\n訂單編號: %s\n區域: %s\n起始日期: %s\n結束日期: %s\n總金額: %s\n----------------------\n訂位者姓名: %s\n電話: %s\n訂位數量: %s\n----------------------\n%s", order_sn, info.Region, info.Start, info.End, info.PaymentTotal, info.UserName, info.PhoneNumber, info.Amount, remit)

	}
	fmt.Println("Relpy_message", reply_mes)

	bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply_mes)).Do()

}

func reply_User_All_Orders(bot *linebot.Client, event *linebot.Event) {
	fmt.Println("reply_User_All_Orders")

	orders, _ := order.GetOrdersByUserID(event.Source.UserID)

	for _, r := range orders {
		fmt.Println(r)

	}
	if len(orders) == 0 {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("您尚未有訂單記錄唷！如有喜歡的營位，請儘速訂位")).Do()
	} else {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("My Orders",
			&linebot.CarouselTemplate{
				Columns:          Carousel_Orders(orders),
				ImageAspectRatio: "rectangle",
				ImageSize:        "cover",
			})).Do()
	}
}

func Carousel_Orders(orders []order.Order) (c_t []*linebot.CarouselColumn) {

	for _, o := range orders {
		start := o.Checkin.Format("2006-01-02")
		end := o.Checkout.Format("2006-01-02")
		camp, _ := product.GetById(int64(o.ProductId))
		var remit string
		var status_mes string
		if o.BankConfirmStatus == order.BankStatus_Unreport {
			status_mes = fmt.Sprintf("狀態:%s (點此回報)", o.BankConfirmStatus)
		} else {
			status_mes = fmt.Sprintf("狀態:%s 後五碼:%s (點此修改)", o.BankConfirmStatus, o.BankLast5Num)

		}
		title := fmt.Sprintf("訂單編號:%s\n區域:%s\n日期:%s~%s\n總金額:%d", o.OrderSN, camp.CampRoundName, start, end, o.PaymentTotal)
		reply_mes := fmt.Sprintf("%s訂位者姓名:%s\n電話:%s\n訂位數量:%d\n%s", title, o.UserName, o.PhoneNumber, o.Amount, remit)
		fmt.Println("reply_mes")
		fmt.Println(reply_mes)
		fmt.Println("status")
		fmt.Println(status_mes)
		tmp := linebot.CarouselColumn{

			ImageBackgroundColor: "#000000",

			Text: reply_mes,
			Actions: []linebot.TemplateAction{
				&linebot.PostbackAction{
					Label:       status_mes,
					Data:        fmt.Sprintf("action=report&data=%s", o.OrderSN),
					InputOption: linebot.InputOptionOpenKeyboard,
					FillInText:  fmt.Sprintf("*回報資訊*\n\n訂單編號:%s\n回報帳號後5碼:", o.OrderSN),
				},
			},
		}
		c_t = append(c_t, &tmp)
	}

	for i, r := range c_t {
		fmt.Println(i, r)
		fmt.Println()
	}
	return c_t
}

func report_Bank_Last_FiveNumbers(bot *linebot.Client, event *linebot.Event, report string) {

	split := strings.Split(report, "\n")
	report_map := make(map[string]string)

	for _, r := range split {
		if strings.Contains(r, ":") {
			arr := strings.Split(r, ":")
			if strings.TrimSpace(arr[1]) != "" {
				report_map[strings.TrimSpace(arr[0])] = strings.TrimSpace(arr[1])
			}
			// if arr[1] != "" {
			// 	info_map[arr[0]] = arr[1]
			// }
		}
	}
	fmt.Println("report map ", report_map)
	sn := report_map["訂單編號"]
	numbers := report_map["回報帳號後5碼"]
	o, _ := order.GetOrderByOrderSN(sn)
	if len(strings.TrimSpace(numbers)) != 5 {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("回報帳號有誤，請重新回報")).Do()

	} else {
		o.BankLast5Num = numbers
		o.BankConfirmStatus = order.BankStatus_UnConfirm
		err := order.UpdateOrder(o)
		if err != nil {
			fmt.Println("Update order bankstatus failed")
		}
		reply_User_All_Orders(bot, event)
	}

}

func get_string_data(str string) string {
	i := strings.Index(str, "=")
	tmp := str[i+1:]
	return tmp
}

func parse_string_to_SearchTime(start, end string) (t Search_Time) {
	t.Start, _ = time.Parse("2006-01-02", start)
	t.End, _ = time.Parse("2006-01-02", end)

	return t
}
