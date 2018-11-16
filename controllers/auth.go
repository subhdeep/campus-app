package controllers

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/subhdeep/campus-app/config"
	"github.com/subhdeep/campus-app/models"
)

type loginCred struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Username string `json:"username"`
}

type unencryptedAuth struct {
	smtp.Auth
}

// Implements unencryptedAuth
func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

var (
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	blockKey []byte
	hashKey  = []byte(config.CookieSecret)
	sc       = securecookie.New(hashKey, blockKey)
	validate = validator.New()
)

// Login is used to perform the login of a user
func Login(ctx iris.Context) {

	user := loginCred{}
	errReq := ctx.ReadJSON(&user)
	if errReq != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	err := validate.Struct(user)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	if !checkLoginCred(&user, ctx) {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}
	username := user.Username
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	expiry := (7 * 24 * time.Hour)
	ctx.SetCookieKV("username", username, iris.CookieEncode(sc.Encode), iris.CookieExpires(expiry))
	ctx.SetCookieKV("timestamp", timestamp, iris.CookieEncode(sc.Encode), iris.CookieExpires(expiry))

	ctx.JSON(loginResponse{
		Username: username,
	})

}

func Logout(ctx iris.Context) {
	ctx.RemoveCookie("username")
	ctx.RemoveCookie("timestamp")
}

// Check is used to check if user is authenticated at the beginning of connection
func Check(ctx iris.Context) {
	username := ctx.Values().Get("userID").(models.Username)
	ctx.JSON(loginResponse{
		Username: string(username),
	})
}

func checkLoginCred(cred *loginCred, ctx iris.Context) bool {
	hostname := config.SMTPHost
	port := config.SMTPPort
	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		ctx.Application().Logger().Warnf("Error occurred while checking login credentials: %s", err.Error())
		return false
	}
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
	return err == nil
}
