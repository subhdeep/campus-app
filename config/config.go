package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// HTTPPort on which the server will run
var HTTPPort int

// CookieSecret which is used for hashing cookies
var CookieSecret string

// SMTPHost for authenticating login
var SMTPHost string

// SMTPPort for authenticating login
var SMTPPort int

// InitConfig initialises the configuration of the app
func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	HTTPPort = viper.GetInt("http.port")
	CookieSecret = viper.GetString("secret.value")

	SMTPPort = viper.GetInt("smtp.port")
	SMTPHost = viper.GetString("smtp.host")
}
