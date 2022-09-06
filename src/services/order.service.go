package services

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type OrderService interface {
	GetOrders() ([]*models.Order, error)
	GetAmbassadorsRevenue(ambassadorID int) (float64, error)
}

type orderService struct {
	orderRepository repositories.OrderRepository
}

func NewOrderService(orderRepository repositories.OrderRepository) OrderService {
	return &orderService{orderRepository: orderRepository}
}

func (s *orderService) GetOrders() ([]*models.Order, error) {
	orders, err := s.orderRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *orderService) GetAmbassadorsRevenue(ambassadorID int) (float64, error) {
	revenue, err := s.orderRepository.CalculateAmbassadorRevenue(ambassadorID)
	if err != nil {
		return 0, nil
	}
	return revenue, nil
}
