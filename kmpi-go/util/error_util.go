package util

import (
	"kmpi-go/log"
	"runtime/debug"
)

func CheckType(i interface{}) string {
	var res string
	switch v := i.(type) {
	case int:
		res = "int type"
	case string:
		return v
	case error:
		res = v.Error()
	default:
		res = "not found type"
	}
	return res
}

/**
返回错误值
*/
func ReturnError() {
	if e := recover(); e != nil {

		prefix := CheckType(e)

		errMessage := prefix + string(debug.Stack())

		log.Error(errMessage, e)

	}

}
