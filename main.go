package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"icarus/database"
	"icarus/router"
	"icarus/user"
)

func init() {
	// return current database
	database.Db.Migrator().CurrentDatabase()

	// Migrate: run auto migration for given models, will only add missing field, won't delete/change current data
	err := database.Db.AutoMigrate(&user.User{}, &user.Department{})
	if err != nil {
		return
	}
	//database.Db.Migrator().RenameTable("users", "user")
	//database.Db.Migrator().RenameTable("departments", "department")
}

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
	err := app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./conf/iris_dev.tml")))
	if err != nil {
		return
	}
}
