package models


import (
	"time"
)

type Quiz struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Date        time.Time      `gorm:"type:date;uniqueIndex;not null" json:"date"` // Quiz for specific date
	Title       string         `gorm:"size:200;not null" json:"title"`             // "Daily Brief - Feb 18, 2026"
	Description string         `gorm:"size:500" json:"description"`
	Difficulty  string         `gorm:"size:20;default:'medium'" json:"difficulty"` // easy, medium, hard
	IsPublished bool           `gorm:"default:false" json:"is_published"`          // Cron publishes at specific time
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relationships
	Questions []Question `gorm:"foreignKey:QuizID" json:"questions"`
}

type Question struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	QuizID       uint           `gorm:"not null;index" json:"quiz_id"`
	ArticleID    *uint          `gorm:"index" json:"article_id"`        
	QuestionText string         `gorm:"type:text;not null" json:"question_text"`
	QuestionType string         `gorm:"size:20;default:'mcq'" json:"question_type"` 
	Options      []Option       `gorm:"foreignKey:QuestionID" json:"options"`
	Explanation  string         `gorm:"type:text" json:"explanation"`  
	Points       int            `gorm:"default:1" json:"points"`
	CreatedAt    time.Time      `json:"created_at"`

	// For True/False or MCQ
	CorrectOptionID *uint      `json:"-"` // Set after options are created
}

type Option struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	QuestionID uint      `gorm:"not null;index" json:"question_id"`
	Text       string    `gorm:"size:500;not null" json:"text"`
	IsCorrect  bool      `gorm:"default:false" json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
}