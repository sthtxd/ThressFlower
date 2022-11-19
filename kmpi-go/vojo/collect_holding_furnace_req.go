package vojo

type CollectHoldingFurnaceReq struct {
	HoldingFurnaceId   *string  `form:"holdingFurnaceId" binding:"required"`
	Status             *int     `form:"status" binding:"required"  json:"status"`
	CurrentTemperature *float64 `form:"currentTemperature" binding:"required"`
	TargetTemperature  *float64 `form:"targetTemperature" binding:"required"`
	OperationType      *int     `form:"operationType" binding:"required"`
	MaxTemperature     *float64 `form:"maxTemperature" binding:"required"`
	MinTemperature     *float64 `form:"minTemperature" binding:"required"`
}
