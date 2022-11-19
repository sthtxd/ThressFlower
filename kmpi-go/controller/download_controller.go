package controller

import (
	"github.com/gin-gonic/gin"
	"kmpi-go/log"
	"kmpi-go/service"
	"kmpi-go/vojo"
	"net/http"
)

func CheckDownloadId(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	response, err := service.CheckDownloadId(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("CheckDownloadId error", err.Error())

	} else {
		res.Message = response
	}

	c.JSON(http.StatusOK, res)
}
func DownloadExcel(c *gin.Context) {
	var res vojo.BaseRes
	err := service.DownloadExcel(c)

	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("ShowImage error", err)
		c.JSON(http.StatusOK, res)
	}
}
