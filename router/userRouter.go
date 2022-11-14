package router

import (
	"icarus/user"

	"github.com/kataras/iris/v12/mvc"
)

func UserRouter(app *mvc.Application) {
	//app.Router.Use(user.BasicAuth)
	// Use memory data in datasource to build Users database (Debug mode?)
	repo := user.NewUserRepository()
	userService := user.NewUserService(repo)
	app.Register(
		userService,
	)
	app.Handle(new(user.Controller))
}
