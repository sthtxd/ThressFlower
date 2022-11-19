package dao

type WeightRfidLogDao struct {
	ID        int            `gorm:"primary_key;auto_increment" json:"id"`
	DeviceId  string         `gorm:"type:varchar(50);not null;index:device_id_idx;comment:'炉号'" json:"deviceId"`
	Rfid      int            `gorm:"type:int(10);not null;default:0;comment:'站点'" json:"Rfid"`
	Weight    int            `gorm:"type:int(10);not null;default:0;index:status_idx;comment:'重量'" json:"weight"`
	Timestamp SelfFormatTime `gorm:"type:timestamp(3) on update current_timestamp(3);omitempty;index:timestamp_idx;default:current_timestamp(3);comment:'时间戳'" json:"timestamp"`
}

func (p WeightRfidLogDao) TableName() string {

	return "weight_Rfid_log"
}
