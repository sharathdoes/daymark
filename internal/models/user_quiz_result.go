package models

import (
	"time"
)

type UserQuizResult struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"index" json:"user_id"`
	QuizID         uint      `gorm:"index" json:"quiz_id"`
	Score          int       `json:"score"`
	TotalQuestions int       `json:"total_questions"`
	Difficulty     string    `json:"difficulty"`
	Categories     string    `json:"categories"` // e.g., comma-separated list or JSON
	CreatedAt      time.Time `json:"created_at"`
}
