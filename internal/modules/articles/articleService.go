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



func (s *Service) CreateArticlesOfCategories( categories []string) error {
		ctx:= context.Background()

	if s.feedSourceService == nil {
		return fmt.Errorf("feed source service is not initialized")
	}

	rssFeeds, err := s.feedSourceService.GetFeedSourcesByCategory(ctx, categories)
	if err != nil {
		return err
	}

	articlesList, err := utils.FetchArticlesFromFeeds(rssFeeds, s.repo.LinkExists)
	if err != nil {
		return err
	}

	
	if err := s.repo.CreateArticle(ctx, articlesList); err != nil {
		return err
	}
	
	return nil
}