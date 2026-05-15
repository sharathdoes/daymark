package articles

import (
	"daymark/config"
	"daymark/internal/modules/category"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	svc := NewService(repo)
	catRepo := category.NewRepository(db)
	catSvc := category.NewService(catRepo)
	h := NewHandler(svc, catSvc)
	r.GET("/articles/today", h.GetTodayArticles)
}
