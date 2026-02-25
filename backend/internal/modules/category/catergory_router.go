package category

import (
	"daymark/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	serv := NewService(repo)
	h := NewHandler(serv)

	g := r.Group("/category")
	{
		g.GET("/", h.GetCategories)
		g.GET("/:id", h.GetCategoryByID)
		g.POST("/", h.CreateCategory)
	}
}
