package controllers

import (
	"fmt"
	"net/smtp"

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

// unencryptedAuth struct
type unencryptedAuth struct {
	smtp.Auth
}

// Implements unencryptedAuth
func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

// Login of a user
func Login(ctx iris.Context) {

	user := LoginCred{}
	errReq := ctx.ReadJSON(&user)
	if errReq != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Text(errReq.Error())
	} else {
		if checkLoginCred(&user) {
			fmt.Println("Connection is successfull")
			ctx.StatusCode(iris.StatusOK)
		} else {
			ctx.StatusCode(iris.StatusForbidden)
		}
	}
}

func checkLoginCred(cred *LoginCred) bool {
	hostname := "smtp.cc.iitk.ac.in"
	conn, err := smtp.Dial(hostname + ":25")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Connection Failed")
		return false
	}
	fmt.Println(cred)
	defer conn.Close()
	auth := unencryptedAuth{
		smtp.PlainAuth(
			"",
			cred.Username,
			cred.Password,
			hostname,
		),
	}
	err = conn.Auth(auth)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Authentication Failed")
		return false
	}
	return true
}
