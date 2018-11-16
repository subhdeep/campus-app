package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kataras/iris"

	"github.com/subhdeep/campus-app/config"
	"github.com/subhdeep/campus-app/router"
)

func main() {
	app := iris.Default()

	router.CampusAppRoutes(app)

	// listen and serve on http://0.0.0.0:8080.
	app.Run(iris.Addr(fmt.Sprintf("0.0.0.0:%d", config.HTTPPort)))
}
