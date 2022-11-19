package service

import (
	"fmt"
	"kmpi-go/dao"
	"kmpi-go/vojo"

	"github.com/gin-gonic/gin"
)

func GetTemperature(c *gin.Context) (*vojo.GetTemperature, error) {
	furnaceId := c.Query("furnaceId")
	if furnaceId == "" {

		return nil, fmt.Errorf("furnaceId is null")
	}
	//furnaceTypeInt, err := strconv.Atoi(furnaceType)
	// if err != nil {
	// 	return nil, err
	// }
	MaxTemperature, MinTemperature, err := searchDbByType(furnaceId)
	if err != nil {
		return nil, err
	}
	var res vojo.GetTemperature
	res.MaxTemperature = MaxTemperature
	res.MinTemperature = MinTemperature

	return &res, nil
}
func searchDbByType(furnaceId string) (*int, *int, error) {

	var res dao.AllDeviceDataDao
	err := dao.CronDb.Where("device_id=?", furnaceId).First(&res).Error
	return &res.MaxTemperature, &res.MinTemperature, err

}
func SaveFurnaceTemperature(c *gin.Context) error {

	var req vojo.SaveFurnaceTemReq

	err := c.BindJSON(&req)
	if err != nil {
		return err
	}

	var tableHolding dao.AllDeviceDataDao
	data := map[string]interface{}{
		"max_temperature": *req.MaxTemperature, "min_temperature": *req.MinTemperature}

	// realData := dao.HoldingFurnaceDao{
	// 	MaxTemperature:    *req.MaxTemperature,
	// 	MinTemperature:    *req.MinTemperature,
	// 	TargetTemperature: *req.TargetTemperature,
	// }
	err = dao.CronDb.Model(&tableHolding).Where("device_id=?", req.DeviceId).Updates(data).Error

	if err != nil {
		return err
	}
	furnaceLog := &dao.DeviceLogDao{
		DeviceId:       *req.DeviceId,
		MaxTemperature: *req.MaxTemperature,
		MinTemperature: *req.MinTemperature,
	}

	return dao.CronDb.Create(furnaceLog).Error

}
