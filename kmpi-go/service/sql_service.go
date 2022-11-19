package service

import (
	"kmpi-go/dao"
	"strconv"
)

//func InsertSqlserver(c *gin.Context) error {
func InsertSqlserver(sName string, tempVisible int, curTemp int, weightVisible int, nWeight int, hydrogenVisible int, hydrogen float64, o2Visible int, o2 float64, nId int) error {
	//guid := uuid.NewV4().String()
	//updateSqlserver()

	// furnaceId := strconv.Itoa(ntype)
	// furnaceLog := &dao.FurnaceTestDao{
	// 	FurnaceId:          furnaceId,
	// 	Weight:             int(nWeight),
	// 	CurrentTemperature: float64(nTemp),
	// }

	// return dao.CronDb.Create(furnaceLog).Error

	data := make(map[string]interface{})
	data["device_name"] = sName
	data["temperature_visible"] = tempVisible
	data["current_temperature"] = curTemp
	data["weight_visible"] = weightVisible
	data["weight"] = nWeight
	data["hydrogen_visible"] = hydrogenVisible
	data["hydrogen"] = hydrogen
	data["oxygen_visible"] = o2Visible
	data["oxygen"] = o2

	result := dao.CronDb.Model(dao.AllDeviceDataDao{}).Where("device_id=?", nId).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	deviceLog := &dao.DeviceLogDao{
		DeviceId:           strconv.Itoa(nId),
		Weight:             nWeight,
		CurrentTemperature: float64(curTemp),
		Oxygen:             o2,
		Hydrogen:           hydrogen,
	}
	return dao.CronDb.Create(deviceLog).Error
	//dao.CronDb.Create(holdingFurnace)
	//err := dao.CronDb.Exec("replace into holding_furnace (furnace_id,status,current_temperature) VALUES (?,?,?)", 1, 2, 3).Error

}

func InsertRfid(nRfid int, nId int) error {
	data := make(map[string]interface{})

	data["rfid"] = nRfid
	result := dao.CronDb.Model(dao.WeightRfidDao{}).Where("device_id=?", nId).Updates(data)
	if result != nil {
		return result.Error
	}

	sensorLog := &dao.WeightRfidLogDao{
		DeviceId: strconv.Itoa(nId),

		Rfid: nRfid,
	}
	return dao.CronDb.Create(sensorLog).Error
}

func InsertWeightSensor(nWeight int, nId int) error {
	data := make(map[string]interface{})
	data["weight"] = nWeight
	result := dao.CronDb.Model(dao.WeightRfidDao{}).Where("device_id=?", nId).Updates(data)
	if result != nil {
		return result.Error
	}

	sensorLog := &dao.WeightRfidLogDao{
		DeviceId: strconv.Itoa(nId),
		Weight:   nWeight,
	}
	return dao.CronDb.Create(sensorLog).Error
}

// func updateSqlserverNew(ntype int, nWeight uint16, nTemp uint16) error {
// 	guid := strconv.Itoa(ntype)

// 	var tableHolding dao.FurnaceTestDao
// 	data := map[string]interface{}{
// 		"max_temperature": nWeight, "min_temperature": nTemp, "target_temperature": nTemp, "timestamp": time.Now()}

// 	err := dao.CronDb.Model(&tableHolding).Where("furnace_id=?", guid).Updates(data).Error

// 	return err
// }
