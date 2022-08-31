package database

import (
	"log"

	"github.com/yogasab/go-monolith-ambassador/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:root@tcp(db:3306)/ambassador"))
	if err != nil {
		panic(err)
	}
	log.Println("Database connected successfully")
}

func AutoMigrate() {
	DB.AutoMigrate(models.User{})
}
