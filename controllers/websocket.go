package controllers

import (
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/websocket"
)

var ws *websocket.Server

func init() {
	ws := websocket.New(websocket.Config{})

	ws.OnConnection(func(c websocket.Connection) {
		ctx := c.Context()
		ctx.Application().Logger().Infof("Websocket connection from: %s", ctx.Values().Get("userID").(string))
	})
}

func Websocket() context.Handler {
	return ws.Handler()
}
