package router

import (
	"icarus/user"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

const RoutePrefix = "/v1/api"

func RootRouter(app *iris.Application) {
	userParty := app.Party(RoutePrefix + "/user")
	// userParty.Use(user.AuthenticatedHandler)
	mvc.Configure(userParty, UserRouter)
	app.UseRouter(iris.NewConditionalHandler(isNotTargetPath, user.AuthenticatedHandler))
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
