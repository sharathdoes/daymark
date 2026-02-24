package database

import (
	"daymark/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&models.Category{},
		&models.FeedSource{},
		&models.Article{},
	)
	return db, nil
}
