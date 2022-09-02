package main

import (
	"math/rand"

	"github.com/bxcodec/faker/v3"
	"github.com/yogasab/go-monolith-ambassador/src/database"
	"github.com/yogasab/go-monolith-ambassador/src/models"
)

func main() {
	database.Connect()
	for i := 0; i < 30; i++ {
		product := &models.Product{
			Title:       faker.Username(),
			Description: faker.Sentence(),
			Image:       faker.URL(),
			Price:       float64(rand.Intn(90) + 10),
		}
		database.DB.Create(&product)
	}
}
