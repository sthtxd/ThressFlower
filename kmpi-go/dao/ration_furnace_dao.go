package dao

import "time"

type RationFurnaceDao struct {
	ID                 int       `gorm:"primary_key;auto_increment" json:"id"`
	FurnaceId          string    `gorm:"type:varchar(50);not null;uniqueIndex:furnace_id_idx;comment:'炉号'" json:"furnaceId"`
	Weight             int       `gorm:"type:int(10);not null;default:0;comment:'重量'" json:"weight"`
	CurrentTemperature float64   `gorm:"type:float(10,2);not null;default:0.00;comment:'当前温度'" json:"currentTemperature"`
	TargetTemperature  float64   `gorm:"type:float(10,2);not null;default:0.00;comment:'目标温度'" json:"targetTemperature"`
	MinTemperature     float64   `gorm:"type:float(10,2);not null;default:0.00;comment:'最小温度'" json:"minTemperature"`
	MaxTemperature     float64   `gorm:"type:float(10,2);not null;default:0.00;comment:'最大温度'" json:"maxTemperature"`
	Timestamp          time.Time `gorm:"type:timestamp(3) on update current_timestamp(3);omitempty;default:current_timestamp(3);comment:'时间戳'" json:"timestamp"`
}

func (p RationFurnaceDao) TableName() string {

	return "ration_furnace"
}
