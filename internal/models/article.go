package models

import "time"
type Article struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string
	Link        string    `gorm:"uniqueIndex"`
	Source      string
	Category    string    `gorm:"index"`
	PublishedAt time.Time
	Content     string    `gorm:"type:text"`
	Processed   bool      `gorm:"default:false"`
	CreatedAt   time.Time
}
