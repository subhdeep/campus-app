package router

import (
	"github.com/kataras/iris"

	"github.com/subhdeep/campus-app/controllers"
	"github.com/subhdeep/campus-app/middlewares"
)

// CampusAppRoutes routes
func CampusAppRoutes(app *iris.Application) {
	user := app.Party("/user")
	user.Get("/me", middlewares.IsAuthenticated, controllers.Check)
	user.Post("/login", controllers.Login)

	app.Get("/ws", middlewares.IsAuthenticated, controllers.Websocket())
	app.Get("/messages", middlewares.IsAuthenticated, controllers.GetMessages)
	app.Get("/recents", middlewares.IsAuthenticated, controllers.GetRecents)
}
