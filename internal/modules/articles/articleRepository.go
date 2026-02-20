package articles

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
func (r *Repository) CreateArticle(ctx context.Context, art []models.Article) error {
	return r.db.WithContext(ctx).Create(&art).Error
}

func (r *Repository) LinkExists(link string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Article{}).
		Where("link = ?", link).
		Count(&count).Error

	return count > 0, err
}