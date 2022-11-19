package service

import (
	"fmt"
	"kmpi-go/dao"
	"kmpi-go/log"
	"kmpi-go/util"
	"kmpi-go/vojo"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func SearchFurnaceLog(req *vojo.SearchDevicePagenationLogReq) (*vojo.SearchlogRes, error) {
	var res vojo.SearchlogRes

	var logData []dao.DeviceLogDao
	var err error
	var totalCount int64

	if *req.Limit < 0 || *req.Offset < 0 {
		return nil, fmt.Errorf("the limit  and the offset can not be <0")
	}
	if req.DeviceId == "" {
		if *req.Offset == 0 {
			err = dao.CronDb.Model(&dao.DeviceLogDao{}).Where("timestamp BETWEEN ? AND ?", req.StartDate, req.EndDate).Count(&totalCount).Error
		}
		err = dao.CronDb.Where("timestamp BETWEEN ? AND ?", req.StartDate, req.EndDate).Limit(*req.Limit).Offset(*req.Offset).Find(&logData).Error

	} else {
		if *req.Offset == 0 {
			err = dao.CronDb.Model(&dao.DeviceLogDao{}).Where("device_id = ?", req.DeviceId).Where("timestamp BETWEEN ? AND ?", req.StartDate, req.EndDate).Count(&totalCount).Error
		}
		err = dao.CronDb.Where("timestamp BETWEEN ? AND ?", req.StartDate, req.EndDate).Where("device_id = ?", req.DeviceId).Limit(*req.Limit).Offset(*req.Offset).Find(&logData).Error

	}
	res.LogData = logData
	res.TotalCount = int(totalCount)
	return &res, err
}

func SearchFurnaceOperationLog(req *vojo.SearchDeviceLogReq) ([]dao.DeviceLogDao, error) {
	var res []dao.DeviceLogDao
	var err error
	if req.DeviceId == "" {
		err = dao.CronDb.Where("operation_type in (?,?) AND timestamp BETWEEN ? AND ?", dao.OperationCollectHoldingFurnaceSetTem, dao.OperationSetTemperature, req.StartDate, req.EndDate).Find(&res).Error
	} else {
		err = dao.CronDb.Where("operation_type in (?,?) AND timestamp BETWEEN ? AND ?", dao.OperationCollectHoldingFurnaceSetTem, dao.OperationSetTemperature, req.StartDate, req.EndDate).Where("device_id = ?", req.DeviceId).Find(&res).Error
	}
	return res, err
}

func SaveExcel(req *vojo.SearchDeviceLogReq, operation_type int) (*string, error) {
	uuid := uuid.NewV4().String()
	if req.DeviceId == "99" {
		res, err := searchAgvLog(req)

		if err != nil {
			return nil, err
		}

		go func() {
			defer util.ReturnError()
			err := WriteAgvExcel(uuid, res)
			if err != nil {
				log.Error("WriteExcel error", err.Error())
			}
			furnaceLog := &dao.DownloadDao{
				DownloadId: uuid,
			}
			err = dao.CronDb.Create(furnaceLog).Error
			if err != nil {
				log.Error("Create error", err.Error())
			}
		}()
	} else {
		res, err := searchLog(req, operation_type)

		if err != nil {
			return nil, err
		}

		go func() {
			defer util.ReturnError()
			err := WriteExcel(uuid, res)
			if err != nil {
				log.Error("WriteExcel error", err.Error())
			}
			furnaceLog := &dao.DownloadDao{
				DownloadId: uuid,
			}
			err = dao.CronDb.Create(furnaceLog).Error
			if err != nil {
				log.Error("Create error", err.Error())
			}
		}()
	}
	return &uuid, nil
}
func searchLog(req *vojo.SearchDeviceLogReq, operation_type int) ([]dao.DeviceLogDao, error) {
	var res []dao.DeviceLogDao
	var err error
	if operation_type == vojo.SQL_TYPE_LOG {
		if req.DeviceId == "" {
			err = dao.CronDb.Where("timestamp BETWEEN ? AND ?", req.StartDate, req.EndDate).Find(&res).Error
		} else {
			err = dao.CronDb.Where("timestamp BETWEEN ? AND ?", req.StartDate, req.EndDate).Where("device_id = ?", req.DeviceId).Find(&res).Error
		}
	} else if operation_type == vojo.SQL_TYPE_OPERATION_LOG {
		if req.DeviceId == "" {
			err = dao.CronDb.Where("operation_type in (?,?) AND timestamp BETWEEN ? AND ?", dao.OperationCollectHoldingFurnaceSetTem, dao.OperationSetTemperature, req.StartDate, req.EndDate).Find(&res).Error
		} else {
			err = dao.CronDb.Where("operation_type in (?,?) AND timestamp BETWEEN ? AND ?", dao.OperationCollectHoldingFurnaceSetTem, dao.OperationSetTemperature, req.StartDate, req.EndDate).Where("device_id = ?", req.DeviceId).Find(&res).Error
		}
	}
	return res, err

}
func LogCount() (*int64, error) {
	var count int64

	err := dao.CronDb.Model(&dao.DeviceLogDao{}).Count(&count).Error

	if err != nil {
		return nil, err
	}
	return &count, nil

}
func LogDeleteLast(rowcount int) {
	sqlStr := "DELETE FROM furnace_log ORDER BY  id asc limit ?"
	err := dao.CronDb.Exec(sqlStr, rowcount).Error
	if err != nil {
		log.Error("del task error:%v", err.Error())
	}
}
func SearchFurnaceLogByType(c *gin.Context) ([]dao.DeviceLogDao, error) {
	var req vojo.SearchFurnaceLogByType
	err := c.BindJSON(&req)
	if err != nil {
		return nil, err
	}
	var res []dao.DeviceLogDao
	err = dao.CronDb.Where("operation_type=?", req.FurnaceType).Order("id asc").Limit(100).Find(&res).Error
	return res, err
}

func SearchAgvAndWeightLog(c *gin.Context) (*vojo.SearchAgvLogRes, error) {
	var req vojo.AgvAndWeightReq
	err := c.BindJSON(&req)
	if err != nil {
		return nil, err
	}
	var logData []dao.WeightRfidLogDao
	var res vojo.SearchAgvLogRes
	var totalCount int64

	err = dao.CronDb.Model(&dao.WeightRfidLogDao{}).Where("device_id = ?", req.DeviceId).Where("timestamp BETWEEN ? AND ?", req.StartDate, req.EndDate).Count(&totalCount).Error

	err = dao.CronDb.Where("timestamp BETWEEN ? AND ?", req.StartDate, req.EndDate).Where("device_id = ?", req.DeviceId).Limit(0).Offset(0).Find(&logData).Error

	res.LogData = logData
	res.TotalCount = int(totalCount)
	return &res, err
}

func searchAgvLog(req *vojo.SearchDeviceLogReq) ([]dao.WeightRfidLogDao, error) {
	var logData []dao.WeightRfidLogDao
	err := dao.CronDb.Where("timestamp BETWEEN ? AND ?", req.StartDate, req.EndDate).Where("device_id = ?", "1").Limit(0).Offset(0).Find(&logData).Error
	return logData, err
}
