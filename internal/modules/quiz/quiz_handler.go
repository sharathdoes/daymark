package quiz

import (
	"daymark/config"
	"daymark/internal/models"
	"daymark/internal/modules/articles"
	"daymark/internal/modules/feedSource"
	"daymark/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo           *Repository
	FeedService    *feedSource.Service
	ArticleService *articles.Service
	cfg            config.Config
}

func (h *Handler) GetArticles(c *gin.Context, CategoryIds []uint) ([]models.Article, error) {
	articles, err := h.ArticleService.GetTodayArticlesByCategory(c, CategoryIds)
	if err != nil {
		return nil, err
	}
	if len(articles) != 0 {
		return articles, nil
	}

	feedsources, err := h.FeedService.GetFeedSourcesByCategory(c, CategoryIds)
	if err != nil {
		return nil, err
	}

	articless, err := h.ArticleService.SyncFromFeeds(c, feedsources, CategoryIds)
	if err != nil {
		return nil, err
	}
	return articless, nil

}

// GenerateQuiz godoc
// @Summary      Generate quiz
// @Description  Generate a quiz based on articles from the given categories
// @Tags         quiz
// @Accept       json
// @Produce      json
// @Param        payload  body      CreateQuizDTO  true  "Quiz request"
// @Success      200      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]string
// @Router       /quiz/generate [post]
func (h *Handler) GenerateQuiz(c *gin.Context) {
	var body CreateQuizDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	articles, err := h.GetArticles(c, body.CategoryIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	quiz, err := services.GenerateQuiz(body.NumberofQuestions, body.CategoryIDs, h.cfg.GROQ_API_KEY, body.Difficulty, articles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quiz)
}
