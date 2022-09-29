package line

import (
	"fmt"
	"linebot/internal/model/product"
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

				if product, ok := Is_Name_Exist(text_trimspace); ok {
					tmp := Img_Carousel_CampRound_Info(product)
					bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("img carousel",
						&linebot.ImageCarouselTemplate{
							Columns: tmp,
						})).Do()
				}

				switch {

				case text_trimspace == "營地位置":
					bot.ReplyMessage(event.ReplyToken, linebot.NewLocationMessage("小路露營區", "426台中市新社區崑山里食水嵙6-2號", 24.2402679, 120.7943069)).Do()

				case text_trimspace == "營地資訊":
					tmp := Quick_Reply_CampRoundName()
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("選擇分區").WithQuickReplies(&linebot.QuickReplyItems{
						Items: tmp,
					})).Do()

					// case strings.Contains(text_trimspace, "起始日期"):
					// 	bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("date range", linebot.NewButtonsTemplate("", "", "選擇日期",
					// 		linebot.NewDatetimePickerAction("結束日期", "action=search&type=get_end_time", "date", time.Now().Format("2006-01-02"), "", time.Now().Format("2006-01-02")),
					// 	))).Do()
					// case strings.Contains(text_trimspace, "結束日期"):
					// 	value := Search[event.Source.UserID]

					// default:
					// 	bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text_trimspace)).Do()
					// }
				}
			}
		}

		if event.Type == linebot.EventTypePostback {
			data := Parase_postback(event.Postback.Data)
			switch data.Action {
			case "search":
				value, isExist := Search[event.Source.UserID]
				if !isExist {
					Search[event.Source.UserID] = &Search_Time{}
				}

				reply_date := func() {
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
					init := time.Now().Format("2006-01-02")
					Max := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
					start_init = init
					start_min = init
					end_init = init
					end_min = init
					start_Max = Max
					end_Max = Max
					fmt.Println(start_time, start_init, start_min, start_Max, end_Max, end_time, end_init, end_min)
					start_time = "起始日期 "
					end_time = "結束日期 "
					// switch {
					// case !value.Start.IsZero() && value.End.IsZero():
					// 	date := value.Start.Format("2006-01-02")
					// 	start_time = fmt.Sprintf("起始日期 %s", date)
					// 	end_time = "結束日期 "
					// case value.Start.IsZero() && !value.End.IsZero():
					// 	date := value.End.Format("2006-01-02")
					// 	start_time = "起始日期 "
					// 	end_time = fmt.Sprintf("結束日期 %s", date)
					// default:
					// 	start_time = "起始日期 "
					// 	end_time = "結束日期 "

					// }
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

				switch data.Type {

				case "go":
					fmt.Println("go reply date")
					reply_date()
				case "get_start_time":
					date := event.Postback.Params.Date
					// str := fmt.Sprintf("起始日期:%s", date)
					fmt.Println("get start time", date)
					value.Start, _ = time.Parse("2006-01-02", date)

					reply_date()

				case "get_end_time":
					date := event.Postback.Params.Date

					value.End, _ = time.Parse("2006-01-02", date)
					fmt.Println("Start Time", Search[event.Source.UserID].Start)
					fmt.Println("End Time", Search[event.Source.UserID].End)

					reply_date()

				case "start_search":
					if !value.Start.IsZero() && !value.End.IsZero() {
						delete(Search, event.Source.UserID)
						bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Camp Search",
							&linebot.CarouselTemplate{
								Columns: Camp_Search_Remain(*value),

								ImageAspectRatio: "rectangle",
								ImageSize:        "cover",
							})).Do()
					}
				}

			}
			fmt.Println("data", event.Postback.Data)

		}
	}

}

//快速回覆營位分區名稱
func Quick_Reply_CampRoundName() (q_p []*linebot.QuickReplyButton) {
	fmt.Println("Quick_Reply_CampRoundName")
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
	return q_p
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

func Camp_Search_Remain(t Search_Time) (c_t []*linebot.CarouselColumn) {

	camp_searchs := SearchRemainCamp(t)
	fmt.Println("input search time ", t)
	for i, r := range camp_searchs {
		fmt.Println(i, ":", r.Stocks)
	}
	for _, s := range camp_searchs {
		remain_num := fmt.Sprintf("剩餘 %d 帳", s.RemainMinAmount)

		des := fmt.Sprintf("%s ~ %s\n%s", t.Start.Format("2006-01-02"), t.End.Format("2006-01-02"), remain_num)
		fmt.Println(s.Product.ImageUri[0], s.Product.CampRoundName, des)
		tmp := linebot.CarouselColumn{
			ThumbnailImageURL:    s.Product.ImageUri[0],
			ImageBackgroundColor: "#000000",
			Title:                s.Product.CampRoundName,
			Text:                 des,
			Actions: []linebot.TemplateAction{
				&linebot.PostbackAction{
					Label: "我要訂位",
					//Data:  fmt.Sprintf("action=order&item=%d", s.Product.ID),
					Data: "action=order",
				},
			},
		}
		fmt.Println(tmp)
		c_t = append(c_t, &tmp)
	}

	return c_t
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

func Img_Carousel_CampRound_Info(product product.Product) (c_t []*linebot.ImageCarouselColumn) {
	fmt.Println("Img_Carousel_CampRound_Info", product)
	for _, uri := range product.ImageUri {
		fmt.Println("Img_Carousel_CampRound_Info : URI :", uri)
		c1 := linebot.ImageCarouselColumn{
			ImageURL: uri,
			Action: &linebot.PostbackAction{
				Label: product.CampRoundName,
				Data:  "action=click&itemid=0",
			},
		}

		c_t = append(c_t, &c1)
	}
	return c_t
}

type ParseData struct {
	Action string
	Item   int
	Type   string
}

func Parase_postback(data string) (p_d ParseData) {
	str := strings.Split(data, "&")

	for _, p := range str {
		switch {
		case strings.Contains(p, "action"):
			p_d.Action = get_string_data(p)
		case strings.Contains(p, "item"):
			p_d.Item, _ = strconv.Atoi(p)
		case strings.Contains(p, "type"):
			p_d.Type = get_string_data(p)
		}
	}
	return p_d
}

func get_string_data(str string) string {
	i := strings.Index(str, "=")
	tmp := str[i+1:]
	return tmp
}
