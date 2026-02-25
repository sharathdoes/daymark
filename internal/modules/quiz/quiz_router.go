package quiz

import (
	"daymark/config"
	"daymark/internal/modules/articles"
	"daymark/internal/modules/feedSource"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes wires the quiz HTTP endpoints.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// repositories & services
	quizRepo := NewRepository(db)
	feedRepo := feedSource.NewRepository(db)
	feedSvc := feedSource.NewService(feedRepo)
	articleRepo := articles.NewRepository(db)
	articleSvc := articles.NewService(articleRepo)

	h := &Handler{
		repo:           quizRepo,
		FeedService:    feedSvc,
		ArticleService: articleSvc,
		cfg:            *cfg,
	}

	g := r.Group("/quiz")
	{
		g.POST("/generate", h.GenerateQuiz)
	}
}
