package controllers

import (
	"github.com/kataras/iris"
)

// LoginCred struct
type LoginCred struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required"`
	Password string `json:"password" xml:"username" form:"username" validate:"required"`
}

// LoginResponse struct
type LoginResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login of a user
func Login(ctx iris.Context) {

	user := LoginCred{}
	errReq := ctx.ReadJSON(&user)
	if errReq != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Text(errReq.Error())
	} else {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(LoginResponse{Username: user.Username, Password: user.Password})
	}
}
