package main

import (
	"fmt"
	database "icarus/database/mariadb"
	_ "icarus/database/redis"
	"icarus/project"
	"icarus/router"
	"icarus/task"
	"icarus/user"

	"github.com/kataras/iris/v12"
	// icarus_grpc "icarus/message/grpc"
)

func init() {
	// return current database
	database.Db.Migrator().CurrentDatabase()

	// Migrate: run auto migration for given models, will only add missing field, won't delete/change current data
	err := database.Db.AutoMigrate(&user.User{}, &user.Department{}, &project.Project{}, &project.ProjectMember{}, &project.TestPlan{}, &task.Task{}, &task.TaskContent{})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}

func main() {

	// Initialize & Register
	app := router.Initialize()
	router.Router(app)

	// run server
	config := iris.WithConfiguration(iris.TOML("./conf/iris_dev.tml"))
	err := app.Run(iris.Addr(":6180"), config, iris.WithoutInterruptHandler, iris.WithLowercaseRouting, iris.WithPathIntelligence)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
}
