package main

import (
	database "icarus/database/mariadb"
	"icarus/project"
	"log"

	// "icarus/router"
	"icarus/task"
	"icarus/user"
	// database "icarus/database/mariadb"
	// _ "icarus/database/redis"
	// "github.com/kataras/iris/v12"
	// icarus_grpc "icarus/message/grpc"
)

func init() {
	// return current database
	database.Db.Migrator().CurrentDatabase()

	// Migrate: run auto migration for given models, will only add missing field, won't delete/change current data
	err := database.Db.AutoMigrate(&user.User{}, &user.Department{}, &project.Project{}, &project.ProjectMember{}, &project.TestPlan{}, &task.Task{}, &task.TaskContent{})
	if err != nil {
		return
	}
}

func main() {

	// app := router.InitRouter()

	// // Register Controller
	// router.RootRouter(app)

	// // run server
	// err := app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.TOML("./conf/iris_dev.tml")))
	// if err != nil {
	// 	return
	// }

	// icarus_grpc.ExampleClient()
	log.Println("----------------------------------------------------------------")
}
