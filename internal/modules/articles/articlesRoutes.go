package articles

import (
	"daymark/internal/modules/feedSource"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepository(db)
	Frepo:=feedSource.NewRepository(db)
	FServ:=feedSource.NewService(Frepo)
	serv:=Newservice(repo, FServ)
	h:=NewHandler(serv)
	g:=r.Group("/article")
	{
		g.POST("/create",h.CreateArticlesOfCategories)
	}
}