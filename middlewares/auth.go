package middlewares

import (
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/subhdeep/campus-app/config"
)

var (
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	blockKey []byte
	hashKey  = []byte(config.CookieSecret)
	sc       = securecookie.New(hashKey, blockKey)
	validate = validator.New()
)

// LoginAuthCred struct
type loginAuthCred struct {
	Username  string `json:"username" validate:"required"`
	Timestamp string `json:"timestamp" validate:"required"`
}

// IsAuthenticated is used to check if a request is authorized
func IsAuthenticated(ctx iris.Context) {
	loginAuth := loginAuthCred{
		Username:  ctx.GetCookie("username", iris.CookieDecode(sc.Decode)),
		Timestamp: ctx.GetCookie("timestamp", iris.CookieDecode(sc.Decode)),
	}

	if err := validate.Struct(loginAuth); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	ctx.Values().Save("userID", loginAuth.Username, true)
	ctx.Next()
}
