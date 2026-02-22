package quiz

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

func (r *Repository) SaveQuiz(questions []models.Question) error {
	quiz := models.Quiz{}
	if err := r.db.Create(&quiz).Error; err != nil {
		return err
	}

	for i := range questions {
		questions[i].QuizID = quiz.ID
	}

	return r.db.Create(&questions).Error
}