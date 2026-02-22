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

func (r *Repository) HasArticlesTodayOfCategory(ctx context.Context,categoryId uint) ( []models.Article, error) {
	var articles []models.Article
	err:=r.db.WithContext(ctx).Model(models.Article{}).Where(" categoryId = ?",categoryId).Find(&articles).Error
	if err!=nil {
		return nil, err
	}
	return articles, nil
}