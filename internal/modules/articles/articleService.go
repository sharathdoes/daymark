package articles

import (
	"context"
	"daymark/internal/modules/feedSource"
	"daymark/pkg/utils"
	"fmt"
)

type Service struct {
	repo              *Repository
	feedSourceService *feedSource.Service
}

func Newservice(repo *Repository, feedSourceService *feedSource.Service) *Service {
	return &Service{repo: repo, feedSourceService: feedSourceService}
}

func (s *Service) CreateArticlesOfCategories(ctx context.Context, categories []string) error {
	if s.feedSourceService == nil {
		return fmt.Errorf("feed source service is not initialized")
	}

	rssFeeds, err := s.feedSourceService.GetFeedSourcesByCategory(ctx, categories)
	if err != nil {
		return err
	}

	articlesList, err := utils.FetchArticlesFromFeeds(rssFeeds)
	if err != nil {
		return err
	}

	for _, article := range articlesList {
		if err := s.repo.CreateArticle(ctx, article); err != nil {
			return err
		}
	}
	return nil
}