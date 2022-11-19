package vojo

type AgvAndWeightReq struct {
	StartDate *string `form:"startDate" binding:"required"`
	EndDate   *string `form:"endDate" binding:"required"`
	DeviceId  *string `form:"deviceId" binding:"required"`
}
