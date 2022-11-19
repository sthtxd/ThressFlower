package service

import (
	"kmpi-go/dao"
)

func GetAgvPosition() ([]dao.WeightRfidDao, error) {
	//holdingFurnaceList, err := dao.GetAllHoldingFurnace()
	var foods []dao.WeightRfidDao
	err := dao.CronDb.Find(&foods).Error

	return foods, err
}
