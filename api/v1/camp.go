package v1

import (
	"linebot/internal/errmsg"
	"linebot/internal/model"
	"linebot/internal/repository"
	"linebot/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCampInfo(c *gin.Context) {
	var test model.Camp
	// fmt.Println(c.Query("index"))
	// fmt.Println(c.Query("name"))

	// name := c.Query("name")
	// city := c.Query("city")
	// town := c.Query("town")
	// phone := c.Query("phone")

	// test = model.Camp{

	// 	Name:        name,
	// 	City:        city,
	// 	Town:        town,
	// 	PhoneNumber: phone,
	// }
	// fmt.Println(test)

	err := repository.CreateNewCamp(&test)
	if err != nil {
		response.Response(c, errmsg.ERROR)
		return
	}

	response.Response(c, errmsg.SUCCESS)
}

func GetCampInfo(c *gin.Context) {
	var query = "哈囉營地"

	camp, err := repository.QueryCampByCampName(query)
	if err != nil {
		response.Response(c, errmsg.ERROR)
		return
	}

	if camp.CampName == "" {
		response.Response(c, errmsg.ERROR_ACCOUNT_NOT_EXIST)
		return
	}

	c.JSON(http.StatusOK, camp)
}
