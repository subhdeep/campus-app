package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
	const addr = "postgresql://subho@localhost:26257/web_app?sslmode=disable"
	var err error
	db, err = gorm.Open("postgres", addr)

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&ChatMessage{})
	db.AutoMigrate(&RecentMessage{})
}
