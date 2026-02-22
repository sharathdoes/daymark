package feedSource

import (
	"daymark/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db  *gorm.DB, cfg *config.Config){
	repo:=NewRepository(db)
	service:=NewService(repo)
	handler:=NewHandler(service)
	g:=r.Group("/feed")
	{
		g.POST("/create", handler.CreateFeed)
		g.GET("/ofCategories", handler.GetFeedSourcesByCategory)
	}
}