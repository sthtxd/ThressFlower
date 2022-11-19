package controller

import (
	"kmpi-go/log"
	"kmpi-go/service"
	"kmpi-go/vojo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTemperature(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	response, err := service.GetTemperature(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("GetAgvPosition error", err.Error())

	} else {
		res.Message = response
	}

	c.JSON(http.StatusOK, res)
}
func SaveFurnaceTemperature(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	err := service.SaveFurnaceTemperature(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("SaveFurnaceTemperature error", err.Error())
	}

	c.JSON(http.StatusOK, res)
}
