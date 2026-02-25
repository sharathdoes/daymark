package category

import (
	"context"
	"daymark/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}


func(r *Repository) GetAll(ctx context.Context) ( []models.Category, error) {
	var categories  []models.Category
	if err:= r.db.WithContext(ctx).Model(&models.Category{}).Find(&categories).Error; err!=nil {
		return nil, err
	}
	return categories, nil
}

func(r *Repository) GetById(ctx context.Context, id uint) ( models.Category, error ) {
	var category models.Category
	if err:=r.db.WithContext(ctx).Where("id = ?", id).Find(&category).Error; err!=nil {
		return models.Category{}, err
	}
	return category, nil
}

func (r *Repository) Create(ctx context.Context, category models.Category) (models.Category, error) {
	if err := r.db.WithContext(ctx).Create(&category).Error; err != nil {
		return models.Category{}, err
	}
	return category, nil
}