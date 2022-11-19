package service

import (
	"kmpi-go/dao"
)

func GetAllDeviceData() ([]dao.AllDeviceDataDao, error) {
	//holdingFurnaceList, err := dao.GetAllHoldingFurnace()
	var foods []dao.AllDeviceDataDao
	err := dao.CronDb.Find(&foods).Error

	return foods, err
}
