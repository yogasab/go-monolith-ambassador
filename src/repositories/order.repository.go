package repositories

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindAll() ([]*models.Order, error)
	CalculateAmbassadorRevenue(ambassadorID int) (float64, error)
	FindUserOrders(link *models.Link) ([]*models.Order, error)
	Create(order *models.Order) (*models.Order, error)
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

func (r *orderRepository) CalculateAmbassadorRevenue(ambassadorID int) (float64, error) {
	var orders []models.Order
	if err := r.DB.Preload("OrderItems").Find(&orders, &models.Order{
		UserID:   uint(ambassadorID),
		Complete: true,
	}).Error; err != nil {
		return 0, err
	}

	var ambassadorRevenue float64 = 0.0
	for _, o := range orders {
		for _, order := range o.OrderItems {
			ambassadorRevenue += order.AmbassadorRevenue
		}
	}

	return ambassadorRevenue, nil
}

func (r *orderRepository) FindUserOrders(link *models.Link) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.DB.Preload("OrderItems").Find(&orders, &models.Order{
		Code:     link.Code,
		Complete: true,
	}).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) Create(order *models.Order) (*models.Order, error) {
	if err := r.DB.Create(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}
