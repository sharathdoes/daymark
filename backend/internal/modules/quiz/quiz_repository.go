package quiz

import (
	"daymark/internal/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveQuiz(ctx context.Context, quiz models.Quiz) error {
	return r.db.WithContext(ctx).Create(&quiz).Error
}