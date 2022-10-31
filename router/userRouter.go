package router

import (
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"icaru/datasource"
	"icaru/user"
	"time"
)

func UserRouter(app *mvc.Application) {
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})
	//app.Router.Use(user.BasicAuth)
	repo := user.NewUserRepository(datasource.Users)
	userService := user.NewUserService(repo)
	app.Register(
		userService,
		sessManager.Start,
	)
	app.Handle(new(user.Controller))
}
