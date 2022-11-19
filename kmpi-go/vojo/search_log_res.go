package vojo

import "kmpi-go/dao"

type SearchlogRes struct {
	LogData    []dao.DeviceLogDao `json:"logData"`
	TotalCount int                `json:"totalCount"`
}
