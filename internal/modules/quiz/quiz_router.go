package quiz

import (
	"daymark/config"
	"daymark/internal/modules/articles"
	"daymark/internal/modules/feedSource"
	"daymark/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes wires the quiz HTTP endpoints.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config, sched Scheduler) {
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
		Scheduler:      sched,
	}

	// Public routes (no auth)
	g := r.Group("/quiz")
	{
		g.GET("/view/:id", h.GetQuizByID)
		g.GET("/daily", h.GetDailyQuiz) // Quiz of the Day
	}

	// Protected routes (JWT required)
	protected := r.Group("/quiz")
	protected.Use(utils.AuthMiddleware(cfg.JWT_SECRET))
	{
		protected.POST("/generate", h.GenerateQuiz)
		protected.POST("/results", h.SaveResult)
		protected.GET("/results", h.GetResults)
		protected.POST("/daily/trigger", h.TriggerDailyQuiz) // on-demand generation
	}
}

