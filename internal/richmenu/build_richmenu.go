package richmenu

import (
	"fmt"
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func Build_RichMenu() {

	bot, err := linebot.New(

		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	Delete_All(bot)
	aid_A := "richmenu-alias-a"
	aid_B := "richmenu-alias-b"

	if CheckRichMenuAlias_exist(bot, aid_A) {
		DeleteRichMenuAlias(bot, aid_A)
	}
	if CheckRichMenuAlias_exist(bot, aid_B) {
		DeleteRichMenuAlias(bot, aid_B)
	}

	Richmenu_Id_A := CreatRichMenu_A(bot, aid_B)

	img_path_A := "./internal/richmenu/img/img_A.png"

	Upload_Img(bot, Richmenu_Id_A, img_path_A)

	Richmenu_Id_B := CreatRichMenu_B(bot, aid_A)

	img_path_B := "./internal/richmenu/img/linerichmenu.jpeg"

	Upload_Img(bot, Richmenu_Id_B, img_path_B)

	Set_Default(bot, Richmenu_Id_B)
	//Set_Manager_RichMenu(bot, "1234", Richmenu_Id_A)

	//CreateRichMenuAlias(bot, aid_A, Richmenu_Id_A)
	//CreateRichMenuAlias(bot, aid_B, Richmenu_Id_B)

}

func CreatRichMenu_A(bot *linebot.Client, aid string) string {
	richMenu := linebot.RichMenu{
		Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
		Selected:    false,
		Name:        "richmenu-a",
		ChatBarText: "選單A",
		Areas: []linebot.AreaDetail{
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 1200, Height: 235},
				Action: linebot.RichMenuAction{
					Type:            linebot.RichMenuActionTypeRichMenuSwitch,
					RichMenuAliasID: aid,
					Data:            "action=richmenu-changed-to-b",
				},
				// Action: linebot.RichMenuAction{
				// 	Type: linebot.RichMenuActionTypeMessage,
				// 	Text: "切換至B",
				// },
			},

			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 788, Width: 833, Height: 553},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeMessage,
					Text: "好友專屬優惠",
				},
			},

			{
				Bounds: linebot.RichMenuBounds{X: 833, Y: 788, Width: 833, Height: 553},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeMessage,
					Text: "常見問題",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1666, Y: 788, Width: 833, Height: 553},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeMessage,
					Text: "推薦我們給朋友",
				},
			},
		},
	}

	res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		log.Fatal(err)
	}

	return res.RichMenuID
}

func CreatRichMenu_B(bot *linebot.Client, aid string) string {
	richMenu := linebot.RichMenu{
		Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
		Selected:    false,
		Name:        "richmenu-b",
		ChatBarText: "選單B",
		Areas: []linebot.AreaDetail{
			//{
			// 	Bounds: linebot.RichMenuBounds{X: 1251, Y: 0, Width: 1200, Height: 235},
			// 	Action: linebot.RichMenuAction{
			// 		Type:            linebot.RichMenuActionTypeRichMenuSwitch,
			// 		RichMenuAliasID: aid,
			// 		Data:            "action=richmenu-changed-to-a",
			// 	},
			// 	// Action: linebot.RichMenuAction{
			// 	// 	Type: linebot.RichMenuActionTypeMessage,
			// 	// 	Text: "切換至A",
			// 	// },
			// },
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 234, Width: 833, Height: 553},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeMessage,
					Text: "我要訂位",
				},
			},

			{
				Bounds: linebot.RichMenuBounds{X: 833, Y: 235, Width: 833, Height: 553},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeMessage,
					Text: "營地資訊",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1666, Y: 235, Width: 833, Height: 553},
				Action: linebot.RichMenuAction{
					Type:  linebot.RichMenuActionTypeMessage,
					Label: "營地位置",
					Text:  "營地位置",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 788, Width: 833, Height: 553},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeMessage,
					Text: "我的訂單",
				},
			},

			{
				Bounds: linebot.RichMenuBounds{X: 833, Y: 788, Width: 833, Height: 553},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeURI,
					URI:  "https://www.google.com.tw/maps/place/%E5%B0%8F%E8%B7%AF%E9%9C%B2%E7%87%9F%E5%8D%80/@24.2402679,120.7943069,17z/data=!3m1!4b1!4m5!3m4!1s0x34691bf3c0f25c6b:0x79e36fadfb5136c8!8m2!3d24.240263!4d120.796501?hl=zh-TW",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1666, Y: 788, Width: 833, Height: 553},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeURI,
					URI:  "tel:0909990685",
				},
			},
		},
	}

	res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		log.Fatal(err)
	}

	return res.RichMenuID
}

func Upload_Img(bot *linebot.Client, id, path string) {
	if _, err := bot.UploadRichMenuImage(id, path).Do(); err != nil {
		log.Fatal(err)
	}
}

func Set_Default(bot *linebot.Client, id string) {
	if _, err := bot.SetDefaultRichMenu(id).Do(); err != nil {
		log.Fatal(err)
	}
}

func Set_Manager_RichMenu(bot *linebot.Client, user_webhook_id string, rich_menu_id string) {
	bot.LinkUserRichMenu(user_webhook_id, rich_menu_id)
	if _, err := bot.SetDefaultRichMenu(user_webhook_id).Do(); err != nil {
		log.Fatal(err)
	}
}

func Delete_All(bot *linebot.Client) {
	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		log.Fatal(err)
	}
	for _, richMenu := range res {
		if _, err := bot.DeleteRichMenu(richMenu.RichMenuID).Do(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("delete all done")
}

func CreateRichMenuAlias(bot *linebot.Client, aid, rid string) {
	if _, err := bot.CreateRichMenuAlias(aid, rid).Do(); err != nil {
		log.Fatal(err)
	}

}

func DeleteRichMenuAlias(bot *linebot.Client, aid string) {
	if _, err := bot.DeleteRichMenuAlias(aid).Do(); err != nil {
		log.Fatal("delete alias failed", err)
	}

}

func CheckRichMenuAlias_exist(bot *linebot.Client, aid string) bool {
	res, err := bot.GetRichMenuAliasList().Do()
	if err != nil {
		log.Fatal(err)
	}
	for _, alias := range res {
		if alias.RichMenuAliasID == aid {
			return true
		}
	}

	return false
}
