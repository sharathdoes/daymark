package database

import (
	"daymark/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  url,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&models.Category{},
		&models.FeedSource{},
		&models.Article{},
		&models.User{},
		&models.UserQuizResult{},
		&models.Quiz{},
	)
	return db, nil
}
