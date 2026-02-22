package articles

import (
	"context"
	"daymark/internal/models"
	"daymark/internal/services"
	"log"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateArticles(ctx context.Context, articles []models.Article) error {
	return s.repo.CreateArticle(ctx, articles)
}

func (s *Service) GetTodayArticlesByCategory(ctx context.Context, categoryIDs []uint) ([]models.Article, error) {
	return s.repo.HasArticlesTodayOfCategory(ctx, categoryIDs)
}

func (s *Service) SyncFromFeeds(ctx context.Context, feedSources []models.FeedSource, categoryIDs []uint) error {
    fetched := services.FetchArticlesFromFeeds(feedSources)
    if len(fetched) == 0 {
        log.Println("No articles fetched from RSS feeds")
        return nil
    }
    return s.repo.BulkUpsert(ctx, fetched)
}

func (s *Service) GetReadyArticles(ctx context.Context, categoryIDs []uint) ([]models.Article, error) {
    return s.repo.GetReadyArticles(ctx, categoryIDs)
}