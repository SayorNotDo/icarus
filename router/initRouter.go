package router

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func InitRouter() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Logger().SetLevel("Debug")

	// loading view templates path
	app.RegisterView(iris.HTML("./templates", ".html"))
	return app
}
