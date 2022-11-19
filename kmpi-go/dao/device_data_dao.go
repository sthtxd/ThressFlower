package dao

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const DB_TIME_FORMAT = "2006-01-02 15:04:05.000"

type AllDeviceDataDao struct {
	ID                 int            `gorm:"primary_key;auto_increment" json:"id"`
	DeviceId           string         `gorm:"type:int(10);not null;default:0;uniqueIndex:device_id_idx;comment:'设备编号'" json:"deviceId"`
	DeviceName         string         `gorm:"type:varchar(50);not null;uniqueIndex:device_name_idx;comment:'设备名称'" json:"deviceName"`
	DeviceStatus       int            `gorm:"type:int(10);not null;default:0;index:device_status_idx;comment:'设备状态：0,正常状态,1,离线状态'" json:"deviceStatus"`
	TemperatureVisible int            `gorm:"type:int(10);not null;default:0;comment:'温度是否显示：0,不显示,1,显示'" json:"temperatureVisible"`
	CurrentTemperature int            `gorm:"type:int(10);not null;default:0;comment:'当前温度'" json:"currentTemperature" `
	MaxTemperature     int            `gorm:"type:int(10);not null;default:0;comment:'上限温度'" json:"maxTemperature" `
	MinTemperature     int            `gorm:"type:int(10);not null;default:0;comment:'下限温度'" json:"minTemperature" `
	WeightVisible      int            `gorm:"type:int(10);not null;default:0;comment:'重量是否显示：0,不显示,1,显示'" json:"weightVisible"`
	Weight             int            `gorm:"type:int(10);not null;default:0;comment:'重量'" json:"weight"`
	HydrogenVisible    int            `gorm:"type:int(10);not null;default:0;comment:'氢气含量是否显示：0,不显示,1,显示'" json:"hydrogenVisible"`
	Hydrogen           float64        `gorm:"type:float(10,4);not null;default:0.00;comment:'氢气含量'" json:"hydrogen"`
	OxygenVisible      int            `gorm:"type:int(10);not null;default:0;comment:'氧气含量是否显示：0,不显示,1,显示'" json:"oxygenVisible"`
	Oxygen             float64        `gorm:"type:float(10,4);not null;default:0.00;comment:'氧气含量'" json:"oxygen"`
	Timestamp          SelfFormatTime `gorm:"type:timestamp(3) on update current_timestamp(3);omitempty;default:current_timestamp(3);comment:'时间戳'" json:"timestamp"`
}
type SelfFormatTime struct {
	time.Time
}

func (p AllDeviceDataDao) TableName() string {

	return "holding_furnace"
}

func (t SelfFormatTime) String() string {
	return t.Format(DB_TIME_FORMAT)
}

//重写 MarshaJSON 方法，在此方法中实现自定义格式的转换；程序中解析到JSON
func (t SelfFormatTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(DB_TIME_FORMAT))
	return []byte(formatted), nil
}

//JSON中解析到程序中
func (t *SelfFormatTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+DB_TIME_FORMAT+`"`, string(data), time.Local)
	*t = SelfFormatTime{Time: now}
	return
}

//写入数据库时会调用该方法将自定义时间类型转换并写入数据库
func (t SelfFormatTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

//读取数据库时会调用该方法将时间数据转换成自定义时间类型
func (t *SelfFormatTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = SelfFormatTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
