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
