package dao

const (
	OperationCollectHoldingFurnaceDefault = 1
	OperationCollectRationFurnace         = 0
	OperationSetTemperature               = 2
	OperationCollectHoldingFurnaceSetTem  = 3

	DEVICE_ONLINE  = 0
	DEVICE_OFFLINE = 1
)

type DeviceLogDao struct {
	ID                 int            `gorm:"primary_key;auto_increment" json:"id"`
	DeviceId           string         `gorm:"type:varchar(50);not null;index:device_id_idx;comment:'设备号'" json:"deviceId"`
	Weight             int            `gorm:"type:int(10);not null;default:0;comment:'炉子重量'" json:"weight"`
	Status             int            `gorm:"type:int(10);not null;default:0;index:status_idx;comment:'状态：0,正常状态,1,缺铝状态'" json:"status"`
	CurrentTemperature float64        `gorm:"type:float(10,2);not null;default:0.00;comment:'当前温度'" json:"currentTemperature"`
	MinTemperature     float64        `gorm:"type:float(10,2);not null;default:0.00;comment:'最小温度'" json:"minTemperature"`
	MaxTemperature     float64        `gorm:"type:float(10,2);not null;default:0.00;comment:'最大温度'" json:"maxTemperature"`
	Oxygen             float64        `gorm:"type:float(10,4);not null;default:0;comment:'氧含量'" json:"oxygen"`
	Hydrogen           float64        `gorm:"type:float(10,4);not null;default:0;comment:'氢含量'" json:"hydrogen"`
	DeviceStatus       int            `gorm:"type:int(10);not null;default:0;index:device_status_idx;comment:'设备状态：0,正常状态,1,离线状态'" json:"deviceStatus"`
	Timestamp          SelfFormatTime `gorm:"type:timestamp(3) on update current_timestamp(3);omitempty;index:timestamp_idx;default:current_timestamp(3);comment:'时间戳'" json:"timestamp"`
}

func (p DeviceLogDao) TableName() string {

	return "furnace_log"
}
