package feedSource

import (
	"context"
	"daymark/internal/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) NewFeed(ctx context.Context, feedSource *models.FeedSource) error {
	return r.db.WithContext(ctx).Create(feedSource).Error
}

func (r *Repository) UpdateFeed(ctx context.Context, feedSource *models.FeedSource) error {
	return r.db.WithContext(ctx).Model(&models.FeedSource{}).Where("id = ? AND  deleted_at IS NULL", feedSource.ID).Updates(feedSource).Error
}

func (r *Repository) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.db.WithContext(ctx).
		Model(&models.FeedSource{}).
		Select("TRIM(UNNEST(category))").
		Where("deleted_at IS NULL").
		Where("category IS NOT NULL").
		Distinct().
		Order("TRIM(UNNEST(category)) ASC").
		Pluck("TRIM(UNNEST(category))", &categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) GetFeedSourcesByCategory(ctx context.Context, categories []string) ([]models.FeedSource, error){
	var feedSources []models.FeedSource
	err:=r.db.WithContext(ctx).Model(&models.FeedSource{}).Where("category && ?",pq.Array(categories)).Find(&feedSources).Error
	if err != nil {
		return nil, err
	}
	return feedSources,nil

}