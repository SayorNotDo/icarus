package utils

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func RestfulResponse(code int16, message string, data interface{}) mvc.Result {
	return mvc.Response{
		ContentType: "application/json",
		Object: iris.Map{
			"code":    code,
			"message": message,
			"data":    data,
		},
	}
}
