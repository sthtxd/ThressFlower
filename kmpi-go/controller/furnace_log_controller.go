package controller

import (
	"kmpi-go/log"
	"kmpi-go/service"
	"kmpi-go/vojo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchFurnaceLog(c *gin.Context) {

	var req vojo.SearchDevicePagenationLogReq

	err := c.BindJSON(&req)

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS
	var logList interface{}
	if err == nil {
		logList, err = service.SearchFurnaceLog(&req)
	}
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("searchfurnaceLog error", err.Error())

	} else {
		res.Message = logList
	}
	c.JSON(http.StatusOK, res)

}
func SearchFurnaceOperationLog(c *gin.Context) {

	var req vojo.SearchDeviceLogReq

	err := c.BindJSON(&req)

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS
	var logList interface{}
	if err == nil {
		logList, err = service.SearchFurnaceOperationLog(&req)
	}
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("searchfurnaceLog error", err.Error())

	} else {
		res.Message = logList
	}
	c.JSON(http.StatusOK, res)

}
func SaveOperationExcel(c *gin.Context) {

	var req vojo.SearchDeviceLogReq

	err := c.BindJSON(&req)

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS
	var logList interface{}
	if err == nil {
		logList, err = service.SaveExcel(&req, vojo.SQL_TYPE_OPERATION_LOG)
	}
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("saveExcel error", err.Error())

	} else {
		res.Message = logList
	}
	c.JSON(http.StatusOK, res)

}
func SaveExcel(c *gin.Context) {

	var req vojo.SearchDeviceLogReq

	err := c.BindJSON(&req)

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS
	var logList interface{}
	if err == nil {
		logList, err = service.SaveExcel(&req, vojo.SQL_TYPE_LOG)
	}
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("saveExcel error", err.Error())

	} else {
		res.Message = logList
	}
	c.JSON(http.StatusOK, res)

}
func SearchFurnaceLogByType(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS
	logList, err := service.SearchFurnaceLogByType(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("saveExcel error", err.Error())

	} else {
		res.Message = logList
	}
	c.JSON(http.StatusOK, res)
}

func SearchAgvAndWeight(c *gin.Context) {
	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS
	logList, err := service.SearchAgvAndWeightLog(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("saveExcel error", err.Error())

	} else {
		res.Message = logList
	}
	c.JSON(http.StatusOK, res)
}
