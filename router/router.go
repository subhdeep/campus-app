package router

import (
	"github.com/kataras/iris"
	"github.com/subhdeep/campus-app/controllers"
)

// CampusAppRoutes routes
func CampusAppRoutes(app *iris.Application) {
	user := app.Party("/user")
	user.Post("/login", controllers.Login)
}
