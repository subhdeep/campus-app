package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kataras/iris"
	"github.com/subhdeep/campus-app/models"
)

func GetMessages(ctx iris.Context) {
	user := ctx.Values().Get("userID").(string)
	offsetParam := ctx.URLParam("offset")
	limitParam := ctx.URLParam("limit")

	var err error

	var limit = 10
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			fmt.Printf("error parsing the limit %s", limitParam)
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}
	}

	var offset = time.Now()

	if offsetParam != "" {
		offset, err = time.Parse(time.RFC3339Nano, offsetParam)
		if err != nil {
			fmt.Printf("error parsing the time %v", offsetParam)
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}
	}

	ctx.JSON(models.GetMessages(user, offset, limit))
}
