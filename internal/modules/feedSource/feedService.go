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

func (s *Service) CreateFeed(ctx context.Context, name string, url string, categoryID uint) error {
	feedSource := &models.FeedSource{Name: name, URL: url, CategoryId: categoryID}
	if err := s.repo.NewFeed(ctx, feedSource); err != nil {
		log.Printf("Error creating feed: %v", err)
		return err
	}
	return nil
}

func (s *Service) GetFeedSourcesByCategory(ctx context.Context, CategoryIDs []uint) ([]models.FeedSource, error) {
	return s.repo.GetFeedSourcesByCategory(ctx, CategoryIDs)
}
