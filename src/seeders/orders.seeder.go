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
		var ordersItems []models.OrderItem
		for j := 0; j < rand.Intn(5); j++ {
			price := float64(rand.Intn(90) + 10)
			qty := uint(rand.Intn(5))
			ordersItems = append(ordersItems, models.OrderItem{
				ProductTitle:      faker.Word(),
				Price:             price,
				Quantity:          qty,
				AdminRevenue:      0.9 * price * float64(qty),
				AmbassadorRevenue: 0.1 * price * float64(qty),
			})
		}
		database.DB.Create(&models.Order{
			UserID:          uint(rand.Intn(30) + 1),
			Code:            faker.Username(),
			AmbassadorEmail: faker.Email(),
			FirstName:       faker.FirstName(),
			LastName:        faker.LastName(),
			Email:           faker.Email(),
			Complete:        true,
			OrderItems:      ordersItems,
		})
	}
}
