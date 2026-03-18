package quiz

import (
	"daymark/internal/models"
	"fmt"
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

// GetDailyQuizByDate returns the Quiz assigned as the daily quiz for the given date ("YYYY-MM-DD").
func (r *Repository) GetDailyQuizByDate(ctx context.Context, date string) (*models.Quiz, error) {
	var dq models.DailyQuiz
	if err := r.db.WithContext(ctx).Where("date = ?", date).First(&dq).Error; err != nil {
		return nil, err
	}
	return r.GetQuizByID(ctx, fmt.Sprintf("%d", dq.QuizID))
}

// SaveDailyQuiz persists a DailyQuiz record (one per date). Returns an error if the date already has one.
func (r *Repository) SaveDailyQuiz(ctx context.Context, dq *models.DailyQuiz) error {
	result := r.db.WithContext(ctx).Create(dq)
	if result.Error != nil {
		log.Printf("[repo] SaveDailyQuiz ERROR: %v", result.Error)
		return result.Error
	}
	log.Printf("[repo] SaveDailyQuiz SUCCESS: date=%s quiz_id=%d", dq.Date, dq.QuizID)
	return nil
}

// HasDailyQuizForDate returns true if a DailyQuiz row already exists for the given date.
func (r *Repository) HasDailyQuizForDate(ctx context.Context, date string) bool {
	var count int64
	r.db.WithContext(ctx).Model(&models.DailyQuiz{}).Where("date = ?", date).Count(&count)
	return count > 0
}