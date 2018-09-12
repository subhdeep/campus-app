package main

import (
	"github.com/kataras/iris"

	"github.com/subhdeep/campus-app/config"
	"github.com/subhdeep/campus-app/router"
)

func main() {
	app := iris.Default()

	config.InitConfig()

	router.CampusAppRoutes(app)

	// listen and serve on http://0.0.0.0:8080.
	app.Run(iris.Addr(":" + string(config.HTTPPort)))
}
