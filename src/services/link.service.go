package services

import (
	"github.com/yogasab/go-monolith-ambassador/src/models"
	"github.com/yogasab/go-monolith-ambassador/src/repositories"
)

type LinkService interface {
	GetUserLinks(UserID int) ([]*models.Link, error)
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
