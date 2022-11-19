package dao

import "time"

type WeightRfidDao struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	DeviceId  string    `gorm:"type:varchar(50);not null;uniqueIndex:device_id_idx;comment:'设备号'" json:"deviceId"`
	Weight    int       `gorm:"type:int;not null;default:0;comment:'重量'" json:"weight"`
	Rfid      int       `gorm:"type:int;not null;default:0;comment:'RFID'" json:"rfid"`
	Timestamp time.Time `gorm:"type:timestamp(3) on update current_timestamp(3);omitempty;default:current_timestamp(3);comment:'时间戳'" json:"timestamp"`
}
