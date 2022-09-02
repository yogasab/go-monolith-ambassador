package main

import (
	"github.com/bxcodec/faker/v3"
	"github.com/yogasab/go-monolith-ambassador/src/database"
	"github.com/yogasab/go-monolith-ambassador/src/models"
)

func main() {
	database.Connect()
	for i := 0; i < 30; i++ {
		ambassador := &models.User{
			FirstName:    faker.FirstName(),
			LastName:     faker.LastName(),
			Email:        faker.Email(),
			IsAmbassador: true,
		}
		hashedPassword := ambassador.HashPassword("password")
		ambassador.Password = hashedPassword
		database.DB.Create(&ambassador)
	}
}
