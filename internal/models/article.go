package models

import "time"

type Article struct {
	ID      uint `gorm:"primaryKey"`
	Title   string
	Link    string `gorm:"uniqueIndex"`
	Source  string
	Content string `gorm:"type:text"`
	Categories   []Category `gorm:"many2many:article_categories;" json:"categories,omitempty"`
	FeedSourceID uint
	PublishedAt  time.Time
}
