package router

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

const RoutePrefix = "/v1/api"

func RootRouter(app *iris.Application) {
	mvc.Configure(app.Party(RoutePrefix+"/user"), UserRouter)
}
