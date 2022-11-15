package main

import (
	"icarus/router"
	"icarus/user"

	database "icarus/database/mariadb"
	_ "icarus/database/redis"

	"github.com/kataras/iris/v12"
)

func init() {
	// return current database
	database.Db.Migrator().CurrentDatabase()

	// Migrate: run auto migration for given models, will only add missing field, won't delete/change current data
	err := database.Db.AutoMigrate(&user.User{}, &user.Department{})
	if err != nil {
		return
	}
}

func main() {

	app := router.InitRouter()

	// Register Controller
	router.RootRouter(app)

	// run server
	err := app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./conf/iris_dev.tml")))
	if err != nil {
		return
	}
}
