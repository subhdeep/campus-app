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
	otherUser := ctx.URLParam("username")
	offsetParam := ctx.URLParam("offset")
	limitParam := ctx.URLParam("limit")

	var err error

	var limit = 20
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

	messages := models.GetMessages(user, otherUser, offset, limit)
	if len(messages) == limit {
		ctx.Header("Link", fmt.Sprintf("%s?username=%s&offset=%s&limit=%d", ctx.Path(), otherUser, messages[len(messages)-1].CreatedAt.Format(time.RFC3339Nano), limit))
	}
	ctx.JSON(messages)
}

func GetRecents(ctx iris.Context) {
	user := ctx.Values().Get("userID").(string)
	offsetParam := ctx.URLParam("offset")
	limitParam := ctx.URLParam("limit")

	var err error

	var limit = 50
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

	messages := models.GetRecents(user, offset, limit)

	var res = make([]models.RecentMessagePayload, 0, len(messages))
	for _, msg := range messages {
		if msg.To == user {
			res = append(res, models.RecentMessagePayload{
				UserID:       msg.From,
				FirstMessage: msg,
			})
		} else {
			res = append(res, models.RecentMessagePayload{
				UserID:       msg.To,
				FirstMessage: msg,
			})
		}
	}
	if len(messages) == limit {
		ctx.Header("Link", fmt.Sprintf("<%s/%s?offset=%s&limit=%d>; rel=\"next\"", ctx.Host(), ctx.Path(), messages[len(messages)-1].CreatedAt.Format(time.RFC3339Nano), limit))
	}
	ctx.JSON(res)
}
