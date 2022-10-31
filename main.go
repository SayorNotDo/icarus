package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
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

	// run server
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./conf/iris_dev.tml")))
}
