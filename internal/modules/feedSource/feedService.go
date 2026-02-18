package feedSource

import (
	"context"
	"daymark/internal/models"
	"log"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) CreateFeed(ctx context.Context, Name string, URL string, Category []string )  error {
	FeedSource := &models.FeedSource{Name: Name, URL: URL, Category: Category}
	if err := s.repo.NewFeed(ctx, FeedSource); err != nil {
		log.Print("Error Creating Feed")
	}
	return nil
}
