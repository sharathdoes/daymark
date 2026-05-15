package subscribe

import (
	"daymark/config"
	"daymark/pkg/email"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	sender := &email.Sender{
		SMTPHost: cfg.SMTP_HOST,
		SMTPPort: cfg.SMTP_PORT,
		Username: cfg.SMTP_USERNAME,
		Password: cfg.SMTP_PASSWORD,
		From:     cfg.SMTP_FROM,
	}
	h := NewHandler(repo, sender)
	r.POST("/subscribe", h.Subscribe)
}
