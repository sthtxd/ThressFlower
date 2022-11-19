package controller

import (
	"github.com/gin-gonic/gin"
	"kmpi-go/log"
	"kmpi-go/service"
	"kmpi-go/vojo"
	"net/http"
)

func Login(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS
	userRes, err := service.Login(c)

	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("login err", err.Error())
	} else {
		res.Message = userRes

	}
	c.JSON(http.StatusOK, res)

}
