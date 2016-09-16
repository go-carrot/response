package response

import (
	"strconv"
)

var errorMapCache map[int]string

func SetErrorMap(errorMap map[int]string) {
	errorMapCache = errorMap
}

type ErrorDetail struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func ErrorDetailText(code int) string {
	if val, ok := errorMapCache[code]; ok {
		return val
	} else {
		panic("ErrorCode '" + strconv.Itoa(code) + "' not registered via SetErrorMap()")
	}
}
