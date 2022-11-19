package vojo

type CollectRationFurnaceReq struct {
	FurnaceId          *string  `form:"furnaceId" binding:"required"`
	Weight             *int     `form:"weight" binding:"required" json:"weight"`
	CurrentTemperature *float64 `form:"currentTemperature" binding:"required"`
}
