package services

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type OrderService interface {
	GetOrders() ([]*models.Order, error)
	GetAmbassadorsRevenue(ambassadorID int) (float64, error)
	GetOrdersRankings() ([]interface{}, error)
	CreateOrder(dto *dto.CreateOrderDTO) (*models.Order, error)
}

type orderService struct {
	orderRepository     repositories.OrderRepository
	userRepository      repositories.UserRepository
	linkRepository      repositories.LinkRepository
	productRepository   repositories.ProductRepository
	orderItemRepository repositories.OrderItemRepository
}

func NewOrderService(
	orderRepository repositories.OrderRepository,
	userRepository repositories.UserRepository,
	linkRepository repositories.LinkRepository,
	productRepository repositories.ProductRepository,
	orderItemRepository repositories.OrderItemRepository,
) OrderService {
	return &orderService{
		orderRepository:     orderRepository,
		userRepository:      userRepository,
		linkRepository:      linkRepository,
		productRepository:   productRepository,
		orderItemRepository: orderItemRepository,
	}
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

func (s *orderService) CreateOrder(dto *dto.CreateOrderDTO) (*models.Order, error) {
	link, err := s.linkRepository.FindByCode(dto.Code)
	if err != nil {
		return nil, err
	}

	order := &models.Order{}
	order.FirstName = dto.FirstName
	order.LastName = dto.LastName
	order.Email = dto.Email
	order.Address = dto.Address
	order.Country = dto.Country
	order.City = dto.City
	order.Zip = dto.ZIP
	order.Code = link.Code
	order.UserID = link.UserID
	order.AmbassadorEmail = link.User.Email
	s.orderRepository.Create(order)

	for _, requestProduct := range dto.Products {
		product, err := s.productRepository.FindByID(requestProduct["product_id"])
		if err != nil {
			return nil, err
		}
		total := product.Price * float64(requestProduct["quantity"])

		item := &models.OrderItem{
			OrderID:           uint(order.ID),
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(requestProduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}
		if _, err = s.orderItemRepository.Create(item); err != nil {
			return nil, err
		}
	}
	return order, nil
}
