package service

import (
	"kmpi-go/config"
	"kmpi-go/dao"
	"kmpi-go/log"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func WriteExcel(fileName string, result []dao.DeviceLogDao) error {

	saveFileName := config.FileSaveDir + "/" + fileName + ".xlsx"
	xlsx := excelize.NewFile()
	// Create a new sheet.
	index := xlsx.NewSheet("Sheet1")
	//保存第一行
	categories := map[string]string{"A1": "操作时间", "B1": "设备号", "C1": "状态", "D1": "当前温度", "E1": "重量", "F1": "氧含量", "G1": "氢含量", "H1": "最大温度", "I1": "最大温度"}
	for k, v := range categories {
		xlsx.SetCellValue("Sheet1", k, v)
	}
	for x := range result {

		realRow := strconv.Itoa(x + 2)
		xlsx.SetCellValue("Sheet1", "A"+realRow, result[x].Timestamp.String())
		xlsx.SetCellValue("Sheet1", "B"+realRow, result[x].DeviceId)
		status := ""
		if result[x].Status == 0 {
			status = "正常"
		} else if result[x].Status == 1 {
			status = "缺铝"
		}
		xlsx.SetCellValue("Sheet1", "C"+realRow, status)
		xlsx.SetCellValue("Sheet1", "D"+realRow, result[x].CurrentTemperature)
		xlsx.SetCellValue("Sheet1", "E"+realRow, result[x].Weight)
		xlsx.SetCellValue("Sheet1", "F"+realRow, result[x].Oxygen)
		xlsx.SetCellValue("Sheet1", "G"+realRow, result[x].Hydrogen)
		xlsx.SetCellValue("Sheet1", "H"+realRow, result[x].MinTemperature)
		xlsx.SetCellValue("Sheet1", "I"+realRow, result[x].MaxTemperature)
	}

	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path.
	err := xlsx.SaveAs(saveFileName)
	if err != nil {
		log.Error("", err.Error())
		return err
	}

	return nil
}
func WriteAgvExcel(fileName string, result []dao.WeightRfidLogDao) error {

	saveFileName := config.FileSaveDir + "/" + fileName + ".xlsx"
	xlsx := excelize.NewFile()
	// Create a new sheet.
	index := xlsx.NewSheet("Sheet1")
	//保存第一行
	categories := map[string]string{"A1": "操作时间", "B1": "设备号", "C1": "当前位置", "D1": "重量"}
	for k, v := range categories {
		xlsx.SetCellValue("Sheet1", k, v)
	}
	for x := range result {

		realRow := strconv.Itoa(x + 2)
		xlsx.SetCellValue("Sheet1", "A"+realRow, result[x].Timestamp.String())
		xlsx.SetCellValue("Sheet1", "B"+realRow, result[x].DeviceId)

		xlsx.SetCellValue("Sheet1", "C"+realRow, result[x].Rfid)
		xlsx.SetCellValue("Sheet1", "D"+realRow, result[x].Weight)
	}

	// Set active sheet of the workbook.
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path.
	err := xlsx.SaveAs(saveFileName)
	if err != nil {
		log.Error("", err.Error())
		return err
	}

	return nil
}
