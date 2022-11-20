package router

import (
	"icarus/user"

	"github.com/kataras/iris/v12/versioning"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

const RoutePrefix = "/api/v1"

func RootRouter(app *iris.Application) {
	userParty := app.Party(RoutePrefix + "/user")
	projectParty := app.Party(RoutePrefix + "/project")
	mvc.Configure(userParty, UserRouter)
	mvc.Configure(projectParty, ProjectRouter)
	app.UseRouter(iris.NewConditionalHandler(isNotTargetPath, user.AuthenticatedHandler))
	app.Use(versioning.FromQuery("version", "1.0.0"))
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
