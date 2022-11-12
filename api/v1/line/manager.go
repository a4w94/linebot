package line

import (
	"fmt"
	"linebot/internal/model/order"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func reply_Unconfirm_Order(bot *linebot.Client, event *linebot.Event) {
	unconfirm_order := order.GetAll_Unconfirm_Order()

	if len(unconfirm_order) == 0 {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("目前沒有需要確認的訂單！")).Do()
	} else {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("確認訂單匯款",
			&linebot.CarouselTemplate{
				Columns:          unconfirm_order_Carousel(unconfirm_order),
				ImageAspectRatio: "rectangle",
				ImageSize:        "cover",
			})).Do()
	}
}

func unconfirm_order_Carousel(orders []order.Order) (c_t []*linebot.CarouselColumn) {

	for _, o := range orders {
		reply_mes := o.Reply_Order_Message()
		reply_mes = fmt.Sprintf("%s\n狀態:%s(後五碼:%s) ", reply_mes, o.ConfirmStatus, o.BankLast5Num)
		tmp := linebot.CarouselColumn{

			ImageBackgroundColor: "#000000",

			Text: reply_mes,
			Actions: []linebot.TemplateAction{
				&linebot.PostbackAction{
					Label: "確認",
					Data:  fmt.Sprintf("action=manager&type=check_order_bank&data=%s", o.OrderSN),
				},
			},
		}

		c_t = append(c_t, &tmp)
	}

	return c_t

}

func (p_d ParseData) check_Unconfirm_Order(bot *linebot.Client, event *linebot.Event) {
	o, err := order.GetOrderByOrderSN(p_d.Data)
	if err != nil {
		log.Println("get order by SN failed")
	}

	if o.ConfirmStatus == order.BankStatus_UnConfirm {
		reply_mes := o.Reply_Order_Message()
		reply_mes = fmt.Sprintf("%s\n狀態:%s(後五碼:%s) ", reply_mes, o.ConfirmStatus, o.BankLast5Num)
		bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("確認該筆訂單匯款", &linebot.ConfirmTemplate{
			Text: reply_mes,
			Actions: []linebot.TemplateAction{
				&linebot.PostbackAction{
					Label: "是",
					Data:  fmt.Sprintf("action=manager&type=check_order_status&data=%s&status=yes", o.OrderSN),
				},
				&linebot.PostbackAction{
					Label: "否",
					Data:  "action=manager&type=check_order_status&status=no",
				},
			},
		})).Do()
	} else if o.ConfirmStatus == order.BankStatus_Confirm {
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("此筆訂單已確認")).Do()
	}
}

func (p_d ParseData) status_Check_Unconfirm_Order(bot *linebot.Client, event *linebot.Event) {

	switch p_d.Status {

	case "no":
		reply_Unconfirm_Order(bot, event)
	case "yes":
		o, err := order.GetOrderByOrderSN(p_d.Data)
		if err != nil {
			log.Println("get order by SN failed")
		}
		o.ConfirmStatus = order.BankStatus_Confirm
		err = order.UpdateOrder(o)
		if err != nil {

			log.Println("update order failed")
			bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("訂單更新狀態失敗")).Do()
		} else {
			bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("訂單確認成功")).Do()

		}

	}
}
