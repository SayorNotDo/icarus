package router

import (
	"context"
	"icarus/exception"
	"icarus/project"
	"icarus/user"
	"log"
	"time"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/versioning"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

const RoutePrefix = "/v1/api"

func Initialize() *iris.Application {
	app := iris.New()
	// CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))
	// app.Use(CSRF)
	app.Use(recover.New())
	app.Use(logger.New())

	app.Logger().SetLevel("Debug")

	app.ConfigureHost(func(su *iris.Supervisor) {
		su.RegisterOnShutdown(func() {
			log.Println("Server terminated")
		})
	})

	// Graceful shutdown
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx)
	})

	return app
}

func Router(app *iris.Application) {
	app.OnErrorCode(iris.StatusNotFound, exception.NotFound)
	app.OnErrorCode(iris.StatusInternalServerError, exception.InternalServerError)
	userParty := app.Party(RoutePrefix + "/user")
	projectParty := app.Party(RoutePrefix + "/project")
	mvc.Configure(userParty, User)
	mvc.Configure(projectParty, Project)
	app.UseRouter(iris.NewConditionalHandler(isNotTargetPath, user.AuthenticatedHandler))
	app.Use(versioning.FromQuery("version", "1.0.0"))
}

func User(app *mvc.Application) {
	//app.Router.Use(user.BasicAuth)
	repo := user.NewUserRepository()
	userService := user.NewUserService(repo)
	app.Register(
		userService,
	)
	app.Handle(new(user.Controller))
}

func Project(app *mvc.Application) {
	repo := project.NewProjectRepository()
	projectService := project.NewProjectService(repo)
	app.Register(
		projectService,
	)
	app.Handle(new(project.Controller))
}

func isNotTargetPath(ctx iris.Context) bool {
	switch ctx.Path() {
	case RoutePrefix + "/user/register":
		return false
	case RoutePrefix + "/user/login":
		return false
	case RoutePrefix + "/user/authenticate":
		return false
	default:
		return true
	}
}
