package vojo

import "kmpi-go/dao"

type SearchAgvLogRes struct {
	LogData    []dao.WeightRfidLogDao `json:"logData"`
	TotalCount int                    `json:"totalCount"`
}
