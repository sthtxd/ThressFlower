package vojo

type SaveFurnaceTemReq struct {
	MaxTemperature *float64 `form:"maxTemperature" binding:"required"`
	MinTemperature *float64 `form:"minTemperature" binding:"required"`
	DeviceId       *string  `form:"deviceId" binding:"required"`
}
