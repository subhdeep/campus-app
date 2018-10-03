package controllers

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/subhdeep/campus-app/config"
)

// LoginCred struct
type LoginCred struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required"`
	Password string `json:"password" xml:"password" form:"password" validate:"required"`
}

// LoginAuthCred struct
type LoginAuthCred struct {
	Username  string `json:"username" validate:"required"`
	Timestamp string `json:"timestamp" validate:"required"`
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

var (
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	blockKey []byte
	hashKey  = []byte(config.CookieSecret)
	sc       = securecookie.New(hashKey, blockKey)
	validate = validator.New()
)

// IsAuthenticated is used to check if a request is authorized
func IsAuthenticated(ctx iris.Context) {
	loginAuth := LoginAuthCred{
		Username:  ctx.GetCookie("username", iris.CookieDecode(sc.Decode)),
		Timestamp: ctx.GetCookie("timestamp", iris.CookieDecode(sc.Decode)),
	}

	if err := validate.Struct(loginAuth); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	ctx.Next()
}

// Login is used to perform the login of a user
func Login(ctx iris.Context) {

	user := LoginCred{}
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
	ctx.StatusCode(iris.StatusOK)

}

func checkLoginCred(cred *LoginCred, ctx iris.Context) bool {
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
