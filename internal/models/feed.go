package models

type FeedSource struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;not null" json:"name"`
	URL  string `gorm:"size:500;not null;unique" json:"url"`
	Categories []Category `gorm:"many2many:feed_source_categories;" json:"categories,omitempty"`
}
