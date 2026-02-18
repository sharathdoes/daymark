package feedSource

import (
	"context"
	"daymark/internal/models"

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

