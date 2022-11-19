package controller

import (
	"kmpi-go/log"
	"kmpi-go/service"
	"kmpi-go/vojo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAgvPosition(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	response, err := service.GetAgvPosition()
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("GetAgvPosition error", err.Error())

	} else {
		res.Message = response
	}

	c.JSON(http.StatusOK, res)
}
