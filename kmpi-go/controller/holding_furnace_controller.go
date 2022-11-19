package controller

import (
	"kmpi-go/log"
	"kmpi-go/service"
	"kmpi-go/vojo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllDeviceData(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	tt, err := service.GetAllDeviceData()
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("GetAllDeviceData error", err.Error())

	} else {
		//data, _ := json.Marshal(tt)
		//log.Info("%s", string(data))
		res.Message = tt
	}

	c.JSON(http.StatusOK, res)

}
