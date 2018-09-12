package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// HTTPPort on which the server will run
var HTTPPort int

// CookieSecret which is used for hashing cookies
var CookieSecret string

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
}
