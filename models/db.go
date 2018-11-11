package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/subhdeep/campus-app/config"
)

var db *gorm.DB

func init() {
	var addr = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.PGUser, "", config.PGHost, config.PGPort, config.PGDB)

	var err error
	db, err = gorm.Open("postgres", addr)

	if err != nil {
		log.Fatalf("Error while connecting to Database: %v", err)
	}

	db.AutoMigrate(&ChatMessage{})
	db.AutoMigrate(&RecentMessage{})
}
