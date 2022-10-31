package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	_ "icaru/database"
	"icaru/router"
)

const RoutePrefix = "/v1/api"

func main() {

	app := router.InitRouter()

	app.Use(recover.New())
	app.Use(logger.New())

	app.Logger().SetLevel("Debug")

	// loading view templates path
	app.RegisterView(iris.HTML("./templates", ".html"))
	// Register Controller
	mvc.Configure(app.Party(RoutePrefix+"/user"), router.UserRouter)

	// run server
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./conf/iris_dev.tml")))
}
