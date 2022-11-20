package router

import (
	"icarus/project"
	"icarus/user"

	"github.com/kataras/iris/v12/mvc"
)

func UserRouter(app *mvc.Application) {
	//app.Router.Use(user.BasicAuth)
	repo := user.NewUserRepository()
	userService := user.NewUserService(repo)
	app.Register(
		userService,
	)
	app.Handle(new(user.Controller))
}

func ProjectRouter(app *mvc.Application) {
	repo := project.NewProjectRepository()
	projectService := project.NewProjectService(repo)
	app.Register(
		projectService,
	)
	app.Handle(new(project.Controller))
}
