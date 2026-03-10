package user

import (
	"daymark/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetByEmail(email string) (*models.User, error) {

	var user models.User

	err := r.db.
		Where("email = ?", email).
		First(&user).Error

	return &user, err
}

func (r *Repository) GetByProvider(provider, providerID string) (*models.User, error) {

	var user models.User

	err := r.db.
		Where("provider = ? AND provider_id = ?", provider, providerID).
		First(&user).Error

	return &user, err
}