package router

import "github.com/kataras/iris/v12"

func InitRouter() *iris.Application {
	app := iris.New()
	return app
}
