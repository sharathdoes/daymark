package articles

import (
	"context"
	"daymark/internal/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *Repository) HasArticlesTodayOfCategory(ctx context.Context,categoryIds []uint) ( []models.Article, error) {
	var articles []models.Article
	now := time.Now()
    startOfDay := time.Date(
        now.Year(),
        now.Month(),
        now.Day(),
        0, 0, 0, 0,
        now.Location(),
    )
	endOfDay := startOfDay.Add(24 * time.Hour)
	err:=r.db.WithContext(ctx).Model(models.Article{}).Where(" category_id  IN ?",categoryIds).Where("published_at >= ? AND published_at <= ?", startOfDay, endOfDay).Find(&articles).Error
	if err!=nil {
		return nil, err
	}
	return articles, nil
}

// BulkUpsert inserts articles & skips duplicates based on unique link
func (r *Repository) BulkUpsert(ctx context.Context, articles []models.Article) ([]models.Article, error) {
	err:= r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "link"}},
			DoNothing: true,	
		}).
		Create(&articles).Error
		return articles, err
}



// GetReadyArticles returns todays articles that have content (ready for quiz)
func (r *Repository) GetReadyArticles(ctx context.Context, categoryIDs []uint) ([]models.Article, error) {
	var articles []models.Article

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.WithContext(ctx).
		Where("category_id IN ?", categoryIDs).
		Where("published_at >= ? AND published_at < ?", startOfDay, endOfDay).
		Where("content != ''").
		Find(&articles).Error

	return articles, err
}