package services

import (
	"github.com/bxcodec/faker/v3"
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/models/dto"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type LinkService interface {
	GetUserLinks(UserID int) ([]*models.Link, error)
	CreateLink(dto *dto.CreateLinkDTO) (*models.Link, error)
}

type linkService struct {
	linkRepository repositories.LinkRepository
}

func NewLinkService(linkRepository repositories.LinkRepository) LinkService {
	return &linkService{linkRepository: linkRepository}
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
