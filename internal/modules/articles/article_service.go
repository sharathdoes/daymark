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
	log.Printf("[articles] CreateArticles count=%d", len(articles))
	return s.repo.CreateArticle(ctx, articles)
}

func (s *Service) GetTodayArticlesByCategory(ctx context.Context, categoryIDs []uint) ([]models.Article, error) {
	log.Printf("[articles] GetTodayArticlesByCategory categories=%v", categoryIDs)
	return s.repo.HasArticlesTodayOfCategory(ctx, categoryIDs)
}

func (s *Service) SyncFromFeeds(ctx context.Context, feedSources []models.FeedSource, categoryIDs []uint) ([]models.Article, error) {
	log.Printf("[articles] SyncFromFeeds sources=%d categories=%v", len(feedSources), categoryIDs)
	fetched, err := services.FetchArticlesFromFeeds(feedSources)
	if err != nil {
		log.Printf("[articles] SyncFromFeeds FetchArticlesFromFeeds error: %v", err)
		return nil, err
	}
	articles, err := s.repo.BulkUpsert(ctx, fetched)
	if err != nil {
		log.Printf("[articles] SyncFromFeeds BulkUpsert error: %v", err)
		return nil, err
	}
	log.Printf("[articles] SyncFromFeeds upserted %d articles", len(articles))
	return articles, nil
}

func (s *Service) GetReadyArticles(ctx context.Context, categoryIDs []uint) ([]models.Article, error) {
	log.Printf("[articles] GetReadyArticles categories=%v", categoryIDs)
	return s.repo.GetReadyArticles(ctx, categoryIDs)
}
