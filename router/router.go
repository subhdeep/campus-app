package router

import (
	"github.com/kataras/iris"

	"github.com/subhdeep/campus-app/controllers"
	"github.com/subhdeep/campus-app/middlewares"
)

// CampusAppRoutes routes
func CampusAppRoutes(app *iris.Application) {
	// or just serve index.html as it is:
	// app.Get("/{f:path}", func(ctx iris.Context) {
	// 	f := ctx.Params().Get("f")
	// 	path := fmt.Sprintf("/home/yash/git/campus-app-frontend/dist/campus-app-frontend/%s", f)
	// 	if _, err := os.Stat(path); !os.IsNotExist(err) {
	// 		ctx.ServeFile(path, false)
	// 	} else {
	// 		ctx.ServeFile("/home/yash/git/campus-app-frontend/dist/campus-app-frontend/index.html", false)
	// 	}
	// })

	api := app.Party("/")
	user := api.Party("/user")
	user.Get("/me", middlewares.IsAuthenticated, controllers.Check)
	user.Post("/login", controllers.Login)
	user.Post("/logout", middlewares.IsAuthenticated, controllers.Logout)

	api.Get("/ws", middlewares.IsAuthenticated, controllers.Websocket())
	api.Get("/messages", middlewares.IsAuthenticated, controllers.GetMessages)
	api.Get("/recents", middlewares.IsAuthenticated, controllers.GetRecents)
	api.Get("/online", middlewares.IsAuthenticated, controllers.IsOnline)
	api.Post("/push-subscription", middlewares.IsAuthenticated, controllers.AddPush)
}
