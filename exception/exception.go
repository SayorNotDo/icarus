package exception

import (
	"github.com/kataras/iris/v12"
)

func NotFound(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"code":    404,
		"message": "Not Found",
		"data":    map[string]string{},
	})
}

func InternalServerError(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"code":    500,
		"message": "Internal Server Error",
		"data":    map[string]string{},
	})
}
