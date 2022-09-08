package services

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type OrderService interface {
	GetOrders() ([]*models.Order, error)
	GetAmbassadorsRevenue(ambassadorID int) (float64, error)
	GetOrdersRankings() ([]interface{}, error)
}

type orderService struct {
	orderRepository repositories.OrderRepository
	userRepository  repositories.UserRepository
}

func NewOrderService(
	orderRepository repositories.OrderRepository,
	userRepository repositories.UserRepository,
) OrderService {
	return &orderService{orderRepository: orderRepository, userRepository: userRepository}
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

func (s *orderService) GetOrdersRankings() ([]interface{}, error) {
	var revenueDetails []interface{}
	var revenueDetail map[string]interface{}

	users, err := s.userRepository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		fullName := user.FirstName + " " + user.LastName
		revenue, err := s.orderRepository.CalculateAmbassadorRevenue(user.ID)
		if err != nil {
			return nil, err
		}
		revenueDetail = map[string]interface{}{"name": fullName, "revenue": revenue}
		revenueDetails = append(revenueDetails, revenueDetail)
	}
	return revenueDetails, nil
}
