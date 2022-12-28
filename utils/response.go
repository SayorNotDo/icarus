package utils

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func Response(code int16, message string, data interface{}) mvc.Result {
	return mvc.Response{
		ContentType: "application/json",
		Object: iris.Map{
			"code":    code,
			"message": message,
			"data":    data,
		},
	}
}

func ResponseError(code int16, message string) mvc.Result {
	return mvc.Response{
		ContentType: "application/json",
		Object: iris.Map{
			"code":    code,
			"message": message,
			"data":    map[string]interface{}{},
		},
	}
}

func ResponseSuccess(message string, data interface{}) mvc.Result {
	return mvc.Response{
		ContentType: "application/json",
		Object: iris.Map{
			"code":    iris.StatusOK,
			"message": message,
			"data":    data,
		},
	}
}
