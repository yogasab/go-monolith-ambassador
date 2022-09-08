package repositories

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	Create(orderItem *models.OrderItem) (*models.OrderItem, error)
}

type orderItemRepository struct {
	DB *gorm.DB
}

func NewOrderItemRepository(DB *gorm.DB) OrderItemRepository {
	return &orderItemRepository{DB: DB}
}

func (r *orderItemRepository) Create(orderItem *models.OrderItem) (*models.OrderItem, error) {
	if err := r.DB.Create(&orderItem).Error; err != nil {
		return nil, err
	}
	return orderItem, nil
}
