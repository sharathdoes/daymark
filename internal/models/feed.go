package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type FeedSource struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	URL         string         `gorm:"size:500;not null;unique" json:"url"`
	Category    pq.StringArray         `gorm:"type:text[]"`
	LastFetched *time.Time     `json:"last_fetched"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Articles []Article `gorm:"foreignKey:FeedSourceID" json:"-"`
}

type Article struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	FeedSourceID uint          `gorm:"not null;index" json:"feed_source_id"`
	Title       string         `gorm:"size:500;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Link        string         `gorm:"size:500;not null" json:"link"`
	body string					
	PublishedAt *time.Time     `json:"published_at"`
	FetchedAt   time.Time      `json:"fetched_at"`
	IsUsed      bool           `gorm:"default:false" json:"is_used"`       
	CreatedAt   time.Time      `json:"created_at"`
	FeedSource FeedSource 		`gorm:"foreignKey:FeedSourceID" json:"-"`
	Questions  []Question `gorm:"foreignKey:ArticleID" json:"-"`
}