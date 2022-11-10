package route

import (
	"fmt"
	"linebot/api/v1/line"
	"linebot/internal/middleware"
	"linebot/internal/richmenu"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	richmenu.Build_RichMenu()

	router := gin.Default()
	// router.Use(gin.BasicAuth(gin.Accounts{
	// 	"admin": "123456",
	// }))
	router.LoadHTMLGlob("web/templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/repeat", line.RepeatHandler)

	router.Any("/callback", line.CampReply)

	router.POST("/po", middleware.JwtMiddleware(), func(c *gin.Context) {
		req, _ := c.Get("email")
		fmt.Println("email", req)
		c.JSON(200, gin.H{"email": req})
	})
	RegisterAccountRoutes(router)
	RegisterOrderRoutes(router)
	// lineroute := router.Group("/callback")
	// lineroute.Any("/", v1.ReplyMessage)

	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.

	return router
}
