package router

import (
	"context"
	"github.com/getsentry/sentry-go"
	sentryiris "github.com/getsentry/sentry-go/iris"
	_ "icarus/docs"
	"icarus/exception"
	"icarus/project"
	"icarus/static"
	"icarus/user"
	"icarus/utils"
	"log"
	"net/http"
	"os"
	"time"

	nriris "github.com/iris-contrib/middleware/newrelic"
	"github.com/iris-contrib/swagger"
	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/versioning"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func Initialize() *iris.Application {
	// sentry handler initialize
	_ = sentry.Init(sentry.ClientOptions{
		Dsn: "",
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					log.Println(req)
				}
			}
			log.Println(event)
			return event
		},
		Debug:            true,
		AttachStacktrace: true,
	})

	relicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Icarus relic"),
		newrelic.ConfigLicense(utils.GetEnv("NEW_RELIC_LICENSE_KEY", "862783809d684c3541dd3bc3cb33fe7e8173NRAL")),
		// newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	app := iris.New()
	//CORS := cors.New(cors.Options{
	//	AllowedOrigins:   []string{"*"},
	//	AllowCredentials: true,
	//})
	//app.Use(CORS)
	// CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))
	// app.Use(CSRF)
	app.Use(recover.New())
	app.Use(logger.New())

	// use sentry handler
	app.Use(sentryiris.New(sentryiris.Options{
		Repanic: true,
	}))
	app.Use(func(ctx iris.Context) {
		if hub := sentryiris.GetHubFromContext(ctx); hub != nil {
			hub.Scope().SetTag("someRandomTag", "maybeYouNeedIt")
		}
		ctx.Next()
	})

	// monitor configuration
	app.Use(nriris.New(relicApp))

	// log configuration
	//logFile, _ := os.Create("icarus-server.log")
	app.Logger().SetLevel("Debug")
	//app.Logger().SetOutput(logFile)

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
		_ = app.Shutdown(ctx)
	})

	return app
}

func Router(app *iris.Application) {
	app.Use(versioning.FromQuery("version", "1.0.0"))
	app.OnErrorCode(iris.StatusNotFound, exception.NotFound)
	app.OnErrorCode(iris.StatusInternalServerError, exception.InternalServerError)
	app.UseGlobal(user.AuthenticatedHandler)
	userParty := app.Party(static.RoutePrefix + "/user")
	projectParty := app.Party(static.RoutePrefix + "/project")
	mvc.Configure(userParty, User)
	mvc.Configure(projectParty, Project)
	removeMiddlewareHandler(app)
}

func index(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"code":    iris.StatusOK,
		"message": "Welcome to Icarus",
		"data":    map[string]string{},
	})
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

func removeMiddlewareHandler(app *iris.Application) {
	config := swagger.Config{
		URL:          "http://localhost:6180/swagger/doc.json",
		DeepLinking:  true,
		DocExpansion: "list",
		DomID:        "#swagger-ui",
		Prefix:       "/swagger",
	}
	swaggerUI := swagger.Handler(swaggerFiles.Handler, config)
	app.Get("/", index).RemoveHandler(user.AuthenticatedHandler)
	app.Get("/swagger", swaggerUI).RemoveHandler(user.AuthenticatedHandler)
	app.Get("/swagger/{any:path}", swaggerUI).RemoveHandler(user.AuthenticatedHandler)
}
