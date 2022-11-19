package vojo

const (
	SQL_TYPE_LOG           = 0
	SQL_TYPE_OPERATION_LOG = 1
)

type SearchDeviceLogReq struct {
	StartDate *string `form:"startDate" binding:"required"`
	EndDate   *string `form:"endDate" binding:"required"`
	DeviceId  string  `form:"deviceId" `
}
type SearchDevicePagenationLogReq struct {
	StartDate *string `form:"startDate" binding:"required"`
	EndDate   *string `form:"endDate" binding:"required"`
	DeviceId  string  `form:"deviceId" `
	Offset    *int    `form:"offset" binding:"required"`
	Limit     *int    `form:"limit"  binding:"required"`
}
