package feedSource

import (
	"context"
	"daymark/internal/models"
	"errors"
	"log"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) CreateFeed(ctx context.Context, name string, url string, categoryIDs []uint) error {
	if len(categoryIDs) == 0 {
		return errors.New("at least one categoryId is required")
	}

	categories, err := s.repo.GetCategoriesByIDs(ctx, categoryIDs)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		return err
	}
	if len(categories) == 0 {
		return errors.New("no valid categories found for given categoryIds")
	}

	feedSource := &models.FeedSource{Name: name, URL: url, Categories: categories}
	if err := s.repo.NewFeed(ctx, feedSource); err != nil {
		log.Printf("Error creating feed: %v", err)
		return err
	}
	return nil
}

func (s *Service) GetFeedSourcesByCategory(ctx context.Context, CategoryIDs []uint) ([]models.FeedSource, error) {
	return s.repo.GetFeedSourcesByCategory(ctx, CategoryIDs)
}
