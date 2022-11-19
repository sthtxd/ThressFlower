package task

import (
	"fmt"
	"kmpi-go/dao"
	"kmpi-go/log"
	"kmpi-go/util"
	"time"

	"github.com/robfig/cron/v3"
)

const OFFLINE_TIME_OUT_MINUTE = 4

func init() {
	c := cron.New(cron.WithSeconds())
	taskCron := "0/10 * * * * ? "
	//taskCron := "0/5 * * * * ? "
	_, err := c.AddFunc(taskCron, func() {
		checkStatus()
	})
	if err != nil {
		log.Error("AddFunc error:", err.Error())
		return
	}

	c.Start()
}
func checkStatus() {
	defer util.ReturnError()

	//	t1 := time.Now() // get current time
	err := searchDataNew()
	if err != nil {
		log.Error("checkStatus error", err)
	}
	// elapsed := time.Since(t1)
	// fmt.Println("App elapsed: ", elapsed)

}

type Result struct {
	FurnaceId string
	Status    int
	Timestamp dao.SelfFormatTime `gorm:"type:timestamp(3) on update current_timestamp(3);omitempty;default:current_timestamp(3);" json:"timestamp"`
}
type ResultNew struct {
	FurnaceId string
	Timestamp dao.SelfFormatTime `gorm:"type:timestamp(3) on update current_timestamp(3);omitempty;default:current_timestamp(3);" json:"timestamp"`
}

func searchDataNew() error {

	var furnaceList []dao.AllDeviceDataDao
	err := dao.CronDb.Where("device_status", dao.DEVICE_ONLINE).Find(&furnaceList).Error
	if err != nil {
		return err
	}
	for _, item := range furnaceList {
		res := ResultNew{}
		sqlSrc := "SELECT device_id,timestamp FROM furnace_log WHERE device_id='%v'   order by timestamp DESC  LIMIT 1"
		sql := fmt.Sprintf(sqlSrc, item.DeviceId)
		err := dao.CronDb.Raw(sql).Find(&res).Error

		if err != nil {
			return err
		}
		timeNow := time.Now()
		timeEnd := res.Timestamp.Time
		currentTimeOut := timeNow.Sub(timeEnd).Minutes()
		if currentTimeOut > OFFLINE_TIME_OUT_MINUTE {
			log.Info("online:设备号：%v，最后一次上传状态时间为：%v", item.DeviceId, res.Timestamp.Time)
			err = updateOffLine(item.DeviceId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// func searchData() error {
// 	list := []Result{}
// 	sqlFormat := "SELECT a.furnace_id,b.status,b.timestamp FROM (SELECT furnace_id FROM holding_furnace where device_status=%v)a LEFT JOIN (SELECT id,STATUS,furnace_id,max(TIMESTAMP)AS timestamp FROM furnace_log WHERE operation_type IN (%v,%v) GROUP BY furnace_id )b ON a.furnace_id=b.furnace_id;"
// 	sql := fmt.Sprintf(sqlFormat, dao.DEVICE_ONLINE, dao.OperationCollectHoldingFurnaceDefault, dao.OperationCollectHoldingFurnaceSetTem)
// 	err := dao.CronDb.Raw(sql).Find(&list).Error
// 	if err != nil {
// 		return err
// 	}
// 	//	log.Info("list is %v", len(list))
// 	for _, item := range list {
// 		timeNow := time.Now()
// 		timeEnd := item.Timestamp.Time
// 		currentTimeOut := timeNow.Sub(timeEnd).Minutes()
// 		//	log.Info("currentTimeOut is %v", currentTimeOut)
// 		if currentTimeOut > OFFLINE_TIME_OUT_MINUTE {
// 			log.Info("设备号：%v，最后一次上传状态时间为：%v", item.FurnaceId, item.Timestamp.Time)
// 			err = updateOffLine(item.FurnaceId)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil

// }

func updateOffLine(furnaceId string) error {
	t1 := time.Now()
	tx := dao.CronDb.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error("updateOffLine error", r)
		}
	}()

	data := make(map[string]interface{})
	data["device_status"] = dao.DEVICE_OFFLINE //1值字段

	//更新设备状态离线
	err := tx.Model(dao.AllDeviceDataDao{}).Where("device_id=?", furnaceId).Updates(data).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	var result dao.DeviceLogDao
	err = tx.Model(dao.DeviceLogDao{}).Where("device_id=?", furnaceId).Order("timestamp desc").Limit(1).First(&result).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if result.DeviceStatus == dao.DEVICE_OFFLINE {
		tx.Commit()
		return nil
	}
	result.DeviceStatus = dao.DEVICE_OFFLINE

	result.ID = 0
	var myDate time.Time

	result.Timestamp.Time = myDate
	err = tx.Model(dao.DeviceLogDao{}).Create(&result).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	elapsed := time.Since(t1)
	log.Info("设备号：%v，更新db的online状态总共花费了：%v", furnaceId, elapsed)

	return nil

}
