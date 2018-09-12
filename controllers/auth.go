package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"time"

	"github.com/kataras/iris"
	"github.com/subhdeep/campus-app/config"
	validator "gopkg.in/go-playground/validator.v9"
)

// LoginCred struct
type LoginCred struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required"`
	Password string `json:"password" xml:"username" form:"username" validate:"required"`
}

// LoginAuthCred struct
type LoginAuthCred struct {
	Username  string `json:"username" validate:"required"`
	Timestamp string `json:"timestamp" validate:"required"`
	Auth      string `json:"auth" validate:"required"`
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

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// IsAuthenticated is used to check if a request is authorized
func IsAuthenticated(ctx iris.Context) {
	loginAuth := LoginAuthCred{
		Username:  ctx.GetCookie("username"),
		Timestamp: ctx.GetCookie("timestamp"),
		Auth:      ctx.GetCookie("auth"),
	}

	if err := validate.Struct(loginAuth); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	if checkHash(loginAuth) {
		ctx.Next()
	} else {
		ctx.StatusCode(iris.StatusForbidden)
	}
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
	if !checkLoginCred(&user) {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}
	username := user.Username
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	// TODO shift to a config file
	secret := config.CookieSecret
	hashValue := []byte(username + ":" + timestamp + ":" + secret)
	hasher := sha256.New()
	hasher.Write(hashValue)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	expiry := (7 * 24 * time.Hour)
	ctx.SetCookieKV("username", user.Username, iris.CookieExpires(expiry))
	ctx.SetCookieKV("timestamp", timestamp, iris.CookieExpires(expiry))
	ctx.SetCookieKV("auth", sha, iris.CookieExpires(expiry))
	ctx.StatusCode(iris.StatusOK)

}

func checkLoginCred(cred *LoginCred) bool {
	hostname := "smtp.cc.iitk.ac.in"
	conn, err := smtp.Dial(hostname + ":25")
	if err != nil {
		fmt.Println(err.Error())
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

func checkHash(loginAuth LoginAuthCred) bool {
	secret := config.CookieSecret
	hashValue := []byte(loginAuth.Username + ":" + loginAuth.Timestamp + ":" + secret)
	hasher := sha256.New()
	hasher.Write(hashValue)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return loginAuth.Auth == sha
}
