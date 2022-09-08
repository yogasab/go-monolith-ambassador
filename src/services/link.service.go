package services

import (
	"log"

	"github.com/bxcodec/faker/v3"
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type LinkService interface {
	GetUserLinks(UserID int) ([]*models.Link, error)
	CreateLink(dto *dto.CreateLinkDTO) (*models.Link, error)
	GetUserLinkStats(UserID int) ([]interface{}, error)
}

type linkService struct {
	linkRepository  repositories.LinkRepository
	orderRepository repositories.OrderRepository
}

func NewLinkService(
	linkRepository repositories.LinkRepository,
	orderRepository repositories.OrderRepository) LinkService {
	return &linkService{linkRepository: linkRepository,
		orderRepository: orderRepository}
}

func (s *linkService) GetUserLinks(UserID int) ([]*models.Link, error) {
	links, err := s.linkRepository.FindAllUser(UserID)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (s *linkService) CreateLink(dto *dto.CreateLinkDTO) (*models.Link, error) {
	link := &models.Link{}
	link.Code = faker.Username()
	link.UserID = uint(dto.UserID)
	for _, productID := range dto.Products {
		product := &models.Product{}
		product.ID = productID
		link.Products = append(link.Products, *product)
	}

	newLink, err := s.linkRepository.Create(link)
	if err != nil {
		return nil, err
	}
	return newLink, nil
}

func (s *linkService) GetUserLinkStats(UserID int) ([]interface{}, error) {
	var results []interface{}
	links, err := s.linkRepository.FindByUserID(UserID)
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		orders, err := s.orderRepository.FindUserOrders(link)
		if err != nil {
			return nil, err
		}
		revenue := 0.0
		for _, order := range orders {
			log.Println(order.OrderItems[0].Price)
			revenue += order.GetTotalPrice()
		}
		results = append(results, map[string]interface{}{
			"link":    link,
			"revenue": revenue,
			"orders":  len(orders),
		})

	}
	return results, nil
}
