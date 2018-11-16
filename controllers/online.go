package controllers

import (
	"github.com/kataras/iris"

	"github.com/subhdeep/campus-app/models"
)

func IsOnline(ctx iris.Context) {
	username := ctx.URLParam("username")
	online, err := models.IsOnline(models.Username(username))

	if err != nil {
		ctx.Application().Logger().Infof("An error occurred for online check: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	if online {
		ctx.StatusCode(iris.StatusOK)
	} else {
		ctx.StatusCode(iris.StatusNotFound)
	}
}
