package models

import (
	"time"
)

type User struct {
	ID uint `gorm:"type:uuid;primaryKey" json:"id"`

	Name  string `json:"name"`
	Email string `gorm:"uniqueIndex" json:"email"`

	PasswordHash string `json:"-"`

	Provider   string `json:"provider"`    // google, github, email
	ProviderID string `json:"provider_id"` // oauth id

	AvatarURL string `json:"avatar_url"`

	EmailVerified     bool      `json:"email_verified"`
	EmailOTP          string    `json:"-"`
	EmailOTPExpiresAt time.Time `json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
