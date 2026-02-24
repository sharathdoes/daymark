package models

import "time"

type Article struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string
	Link        string `gorm:"uniqueIndex"`
	Source      string
	Content     string `gorm:"type:text"`
	CategoryID  uint
	FeedSourceID uint
	PublishedAt time.Time
}

