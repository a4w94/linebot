package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func repeatHandler(r int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello from Go!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

func dbFunc(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error incrementing tick: %q", err))
			return
		}

		rows, err := db.Query("SELECT tick FROM ticks")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			var tick time.Time
			if err := rows.Scan(&tick); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
		}
	}
}

func replyMessage(bot *linebot.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("get query", c.Query("input"))

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
					switch text_trimspace {
					case "hi":
						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("hi 你好")).Do()
					case "sticker":

						bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage("446", "1988")).Do()
					case "img":
						//參數：原始內容網址 預覽圖
						bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage("https://example.com/original.jpg", "https://example.com/preview.jpg")).Do()
					case "video":
						//參數：連結 預覽圖
						bot.ReplyMessage(event.ReplyToken, linebot.NewVideoMessage("https://example.com/original.mp4", "https://example.com/preview.jpg")).Do()
					case "audio":
						//音頻長度（毫秒）
						audiolenth := 6000
						bot.ReplyMessage(event.ReplyToken, linebot.NewAudioMessage("https://example.com/original.m4a", audiolenth)).Do()
					case "button":
						bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("this is a button template", &linebot.ButtonsTemplate{
							ThumbnailImageURL:    "https://example.com/bot/images/image.jpg",
							ImageAspectRatio:     "rectangle",
							ImageSize:            "cover",
							ImageBackgroundColor: "#FFFFFF",
							Title:                "Menu",
							Text:                 "Please select",
							DefaultAction: &linebot.URIAction{
								Label: "View detail",
								URI:   "http://example.com/page/123",
							},
							Actions: []linebot.TemplateAction{
								&linebot.PostbackAction{
									Label: "Buy",
									Data:  "action=buy&itemid=123",
								},
								&linebot.PostbackAction{
									Label: "Add to cart",
									Data:  "action=add&itemid=123",
								},
								&linebot.URIAction{
									Label: "View detail",
									URI:   "http://example.com/page/123",
								},
							},
						})).Do()
					case "confirm":
						bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("confirm template",
							&linebot.ConfirmTemplate{
								Text: "confirm",
								Actions: []linebot.TemplateAction{
									&linebot.MessageAction{
										Label: "Yes",
										Text:  "yes",
									},
									&linebot.MessageAction{
										Label: "No",
										Text:  "no",
									},
								},
							})).Do()
					case "carousel":
						column1 := linebot.CarouselColumn{
							ThumbnailImageURL:    "https://example.com/bot/images/item1.jpg",
							ImageBackgroundColor: "#FFFFFF",
							Title:                "this is menu",
							Text:                 "description",
							DefaultAction: &linebot.URIAction{
								Label: "View detail",
								URI:   "http://example.com/page/123",
							},
							Actions: []linebot.TemplateAction{
								&linebot.PostbackAction{
									Label: "Buy",
									Data:  "action=buy&itemid=111",
								},
								&linebot.PostbackAction{
									Label: "Add to chart",
									Data:  "action=buy&itemid=111",
								},
								&linebot.URIAction{
									Label: "View detail",
									URI:   "http://example.com/page/111",
								},
							},
						}
						column2 := linebot.CarouselColumn{
							ThumbnailImageURL:    "https://example.com/bot/images/item2.jpg",
							ImageBackgroundColor: "#000000",
							Title:                "this is menu",
							Text:                 "description",
							DefaultAction: &linebot.URIAction{
								Label: "View detail",
								URI:   "http://example.com/page/222",
							},
							Actions: []linebot.TemplateAction{
								&linebot.PostbackAction{
									Label: "Buy",
									Data:  "action=buy&itemid=222",
								},
								&linebot.PostbackAction{
									Label: "Add to chart",
									Data:  "action=buy&itemid=222",
								},
								&linebot.URIAction{
									Label: "View detail",
									URI:   "http://example.com/page/222",
								},
							},
						}
						bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("carousel template",
							&linebot.CarouselTemplate{
								Columns:          []*linebot.CarouselColumn{&column1, &column2},
								ImageAspectRatio: "rectangle",
								ImageSize:        "cover",
							})).Do()
					case "img carousel":
						bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("img_carousel",
							&linebot.ImageCarouselTemplate{
								Columns: []*linebot.ImageCarouselColumn{
									{
										ImageURL: "https://example.com/bot/images/item1.jpg",
										Action: &linebot.PostbackAction{
											Label: "Buy",
											Data:  "action=buy&itemid=222",
										},
									},
									{
										ImageURL: "https://example.com/bot/images/item2.jpg",
										Action: &linebot.URIAction{
											Label: "View detail",
											URI:   "http://example.com/page/222",
										},
									},
									{
										ImageURL: "https://example.com/bot/images/item3.jpg",
										Action: &linebot.MessageAction{
											Label: "Yes",
											Text:  "yes",
										},
									},
								},
							})).Do()
					case "flex":
						bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("this is a flex message",
							&linebot.BubbleContainer{
								Type: "bubble",
								Body: &linebot.BoxComponent{
									Type:   "box",
									Layout: "vertical",
									Contents: []linebot.FlexComponent{
										&linebot.TextComponent{
											Type: "text",
											Text: "Hello",
										},
										&linebot.TextComponent{
											Type: "text",
											Text: "World",
										},
									},
								},
							},
						)).Do()
					case "quick reply":

						bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Please Select").WithQuickReplies(&linebot.QuickReplyItems{
							Items: []*linebot.QuickReplyButton{
								{
									ImageURL: "https://example.com/sushi.png",
									Action: &linebot.MessageAction{
										Label: "Sushi",
										Text:  "Sushi",
									},
								},
								{
									ImageURL: "https://example.com/tempura.png",
									Action: &linebot.MessageAction{
										Label: "Tempura",
										Text:  "Tempura",
									},
								},
								{
									Action: &linebot.LocationAction{
										Label: "Send locatiom",
									},
								},
								{
									Action: &linebot.CameraRollAction{
										Label: "Send photo",
									},
								},
								{
									Action: &linebot.CameraAction{
										Label: "Open Camera",
									},
								},
								{
									Action: &linebot.DatetimePickerAction{
										Label: "select Date",
										Data:  "actionId=32",
										Mode:  "date",
									},
								},
							},
						})).Do()

					case "help":
						bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("help template",
							&linebot.CarouselTemplate{
								Columns: []*linebot.CarouselColumn{
									{

										Title: "Test",
										Text:  "options",

										Actions: []linebot.TemplateAction{
											&linebot.MessageAction{
												Label: "hi--文字訊息",
												Text:  "hi",
											},
											&linebot.MessageAction{
												Label: "img--圖像訊息",
												Text:  "img",
											},

											&linebot.MessageAction{
												Label: "video--影片訊息",
												Text:  "video",
											},
										},
									},
									{

										Title: "Test",
										Text:  "options",

										Actions: []linebot.TemplateAction{
											&linebot.MessageAction{
												Label: "audio--音頻訊息",
												Text:  "audio",
											},
											&linebot.MessageAction{
												Label: "button--按鈕模板",
												Text:  "button",
											},
											&linebot.MessageAction{
												Label: "confirm--確認模板",
												Text:  "confirm",
											},
										},
									},
									{

										Title: "Test",
										Text:  "options",

										Actions: []linebot.TemplateAction{
											&linebot.MessageAction{
												Label: "carousel--輪播模板",
												Text:  "carousel",
											},
											&linebot.MessageAction{
												Label: "img carousel--輪播圖像",
												Text:  "img carousel",
											},
											&linebot.MessageAction{
												Label: "flex--渲染模板",
												Text:  "flex",
											},
										},
									},
									{

										Title: "Test",
										Text:  "options",

										Actions: []linebot.TemplateAction{

											&linebot.MessageAction{
												Label: "quick reply--快速回復",
												Text:  "quick reply",
											},
											&linebot.MessageAction{
												// Label: "carousel--輪播模板",
												// Text:  "carousel",
											},
											&linebot.MessageAction{
												// Label: "img carousel--輪播圖像",
												// Text:  "img carousel",
											},
										},
									},
								},
								ImageAspectRatio: "rectangle",
								ImageSize:        "cover",
							},
						)).Do()
					}

				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}

			}
		}
	}
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	} else {
		log.Println("port", port)
	}

	tStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatalf("Error opening database: %q", err)
	// }
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/repeat", repeatHandler(repeat))
	//router.GET("/db", dbFunc(db))
	router.Any("/callback", replyMessage(bot))

	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	router.Run(":" + port)

}
