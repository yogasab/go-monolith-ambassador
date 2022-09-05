package repositories

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindAll() ([]*models.Order, error)
}

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB) OrderRepository {
	return &orderRepository{DB: DB}
}

func (r *orderRepository) FindAll() ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.DB.Preload("OrderItems").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
