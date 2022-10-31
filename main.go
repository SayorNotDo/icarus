package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	_ "icaru/database"
	"icaru/router"
)

func main() {

	app := router.InitRouter()

	app.Use(recover.New())
	app.Use(logger.New())

	app.Logger().SetLevel("Debug")

	// loading view templates path
	app.RegisterView(iris.HTML("./templates", ".html"))
	// Register Controller
	router.RootRouter(app)
	//mvc.Configure(app.Party(RoutePrefix+"/user"), router.UserRouter)
	//mvc.Configure(app.Party(RoutePrefix), router.RootRouter)
	// run server
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./conf/iris_dev.tml")))
}
