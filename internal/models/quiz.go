package models

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Question struct {
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	Answer     int      `json:"answer"`
	ArticleURL string   `json:"article_url"`
}

type CategoryIDs []uint

func (c CategoryIDs) Value() (driver.Value, error) {
	// Store as plain JSON text so the column contains readable JSON, not hex
	bytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (c *CategoryIDs) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		if err := json.Unmarshal(v, c); err != nil {
			log.Printf("[models] CategoryIDs scan failed for bytes: %v raw=%q", err, string(v))
			*c = CategoryIDs{}
			return nil
		}
	case string:
		if err := json.Unmarshal([]byte(v), c); err != nil {
			log.Printf("[models] CategoryIDs scan failed for string: %v raw=%q", err, v)
			*c = CategoryIDs{}
			return nil
		}
	default:
		return fmt.Errorf("failed to scan CategoryIDs")
	}
	return nil
}

type Questions []Question

func (q Questions) Value() (driver.Value, error) {
	// Store as plain JSON text so the column contains readable JSON, not hex
	bytes, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (q *Questions) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		// First try direct JSON
		if err := json.Unmarshal(v, q); err == nil {
			return nil
		}
		// Fallback: value might contain a hex-encoded JSON blob (e.g. "r\\x5b7b22...")
		raw := string(v)
		if idx := strings.Index(raw, "\\x"); idx >= 0 && len(raw) > idx+2 {
			hexPart := strings.TrimSpace(raw[idx+2:])
			if decoded, err := hex.DecodeString(hexPart); err == nil {
				if err := json.Unmarshal(decoded, q); err == nil {
					return nil
				}
			}
		}
		log.Printf("[models] Questions scan failed for bytes (after hex fallback): raw=%q", raw)
		*q = Questions{}
		return nil
	case string:
		// First try direct JSON
		if err := json.Unmarshal([]byte(v), q); err == nil {
			return nil
		}
		// Fallback for hex-encoded JSON stored inside the string (e.g. "r\\x5b7b22...")
		raw := v
		if idx := strings.Index(raw, "\\x"); idx >= 0 && len(raw) > idx+2 {
			hexPart := strings.TrimSpace(raw[idx+2:])
			if decoded, err := hex.DecodeString(hexPart); err == nil {
				if err := json.Unmarshal(decoded, q); err == nil {
					return nil
				}
			}
		}
		log.Printf("[models] Questions scan failed for string (after hex fallback): raw=%q", raw)
		*q = Questions{}
		return nil
	default:
		return fmt.Errorf("failed to scan Questions")
	}
	return nil
}

type Quiz struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Title       string      `json:"title"`
	CategoryIDs CategoryIDs `gorm:"type:text" json:"category_ids"`
	Difficulty  string      `json:"difficulty"`
	Questions   Questions   `gorm:"type:text" json:"questions"`
	CreatedAt   time.Time   `json:"created_at"`
}

// DailyQuiz records which Quiz is selected as the Quiz of the Day.
// Date is stored as "YYYY-MM-DD" and is unique — one quiz per calendar day.
type DailyQuiz struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	QuizID    uint      `gorm:"index" json:"quiz_id"`
	Date      string    `gorm:"uniqueIndex;size:10" json:"date"` // "YYYY-MM-DD"
	CreatedAt time.Time `json:"created_at"`
}
