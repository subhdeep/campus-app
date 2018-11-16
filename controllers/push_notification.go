package controllers

import (
	"io/ioutil"

	"github.com/kataras/iris"

	"github.com/subhdeep/campus-app/models"
)

// AddPush is used to add a push notification
func AddPush(ctx iris.Context) {
	username := ctx.Values().Get("userID").(models.Username)

	rawData, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	models.CreatePushNotification(string(username), string(rawData))
	ctx.StatusCode(iris.StatusAccepted)
}
