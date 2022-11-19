package controller

import (
	"kmpi-go/log"
	"kmpi-go/service"
	"kmpi-go/vojo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminLogin(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	tt, err := service.AdminLogin(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("AdminLogin error", err.Error())
	} else {
		res.Message = tt
	}
	c.JSON(http.StatusOK, res)
}
func DeleteUser(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	err := service.DeleteUser(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("DeleteUser error", err.Error())

	}
	c.JSON(http.StatusOK, res)
}
func AddUser(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	err := service.AddUser(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("AddUser error", err.Error())

	}
	c.JSON(http.StatusOK, res)
}
func UpdateUser(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	err := service.UpdateUser(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("UpdateUser error", err.Error())

	}
	c.JSON(http.StatusOK, res)
}
func GetAllUser(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	users, err := service.GetAllUser(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("GetAllUser error", err.Error())

	} else {
		res.Message = users
	}
	c.JSON(http.StatusOK, res)
}
func ModifyPassword(c *gin.Context) {

	var res vojo.BaseRes
	res.Rescode = vojo.NORMAL_RESPONSE_STATUS

	err := service.ModifyPassword(c)
	if err != nil {
		res.Rescode = vojo.ERROR_RESPONSE_STATUS
		res.Message = err.Error()
		log.Error("GetAllUser error", err.Error())

	}
	c.JSON(http.StatusOK, res)
}
