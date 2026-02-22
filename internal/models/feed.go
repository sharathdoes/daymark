package models

type FeedSource struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	URL         string         `gorm:"size:500;not null;unique" json:"url"`
	CategoryId 	uint
	Category	Category		`gorm:"foreignKey:CategoryId"` 
}
