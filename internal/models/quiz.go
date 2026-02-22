package models

import (
    "time"
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

type Question struct {
    Question   string   `json:"question"`
    Options    []string `json:"options"`
    Answer     string   `json:"answer"`
    ArticleURL string   `json:"article_url"`
}

type Questions []Question

func (q Questions) Value() (driver.Value, error) {
    return json.Marshal(q)
}

func (q *Questions) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("failed to scan Questions")
    }
    return json.Unmarshal(bytes, q)
}

type Quiz struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    CategoryIDs string   `gorm:"type:text" json:"category_ids"` // stored as "1,2,3"
    Difficulty  string   `json:"difficulty"`                     // easy, medium, hard
    Questions   Questions `gorm:"type:jsonb" json:"questions"`   // all 15 questions
    CreatedAt  time.Time `json:"created_at"`
}