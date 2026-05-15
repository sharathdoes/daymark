package models

type Subscriber struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Email       string `gorm:"uniqueIndex" json:"email"`
	CategoryIDs string `gorm:"type:text" json:"category_ids"`
}
