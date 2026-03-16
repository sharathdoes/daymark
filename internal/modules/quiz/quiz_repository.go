package quiz

import (
	"daymark/internal/models"
	"log"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveQuiz(ctx context.Context, quiz *models.Quiz) error {
	log.Printf("[repo] SaveQuiz called: title=%q questions=%d difficulty=%q categoryIDs=%v",
		quiz.Title, len(quiz.Questions), quiz.Difficulty, quiz.CategoryIDs)
	result := r.db.WithContext(ctx).Create(quiz)
	if result.Error != nil {
		log.Printf("[repo] SaveQuiz ERROR: %v", result.Error)
		return result.Error
	}
	log.Printf("[repo] SaveQuiz SUCCESS: inserted quiz.ID=%d", quiz.ID)
	return nil
}

func (r *Repository) GetQuizByID(ctx context.Context, id string) (*models.Quiz, error) {
	var quiz models.Quiz
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&quiz).Error
	if err != nil {
		log.Printf("[repo] GetQuizByID ERROR id=%s err=%v", id, err)
		return nil, err
	}
	log.Printf("[repo] GetQuizByID SUCCESS id=%s title=%q questions=%d", id, quiz.Title, len(quiz.Questions))
	return &quiz, nil
}

func (r *Repository) SaveUserResult(ctx context.Context, result *models.UserQuizResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

func (r *Repository) GetResultsByUserID(ctx context.Context, userID uint) ([]models.UserQuizResult, error) {
	var results []models.UserQuizResult
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&results).Error
	return results, err
}