package quiz

import (
	"daymark/config"
	"daymark/internal/models"
	"daymark/internal/modules/articles"
	"daymark/internal/modules/feedSource"
	"daymark/internal/services"
	"log"
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
	log.Printf("[quiz] GetArticles start - categories=%v", CategoryIds)
	articles, err := h.ArticleService.GetTodayArticlesByCategory(c, CategoryIds)
	if err != nil {
		log.Printf("[quiz] GetArticles error on GetTodayArticlesByCategory: %v", err)
		return nil, err
	}
	if len(articles) != 0 {
		log.Printf("[quiz] GetArticles found %d existing articles", len(articles))
		return articles, nil
	}

	feedsources, err := h.FeedService.GetFeedSourcesByCategory(c, CategoryIds)
	if err != nil {
		log.Printf("[quiz] GetArticles error on GetFeedSourcesByCategory: %v", err)
		return nil, err
	}

	articless, err := h.ArticleService.SyncFromFeeds(c, feedsources, CategoryIds)
	if err != nil {
		log.Printf("[quiz] GetArticles error on SyncFromFeeds: %v", err)
		return nil, err
	}
	log.Printf("[quiz] GetArticles fetched %d articles from feeds", len(articless))
	return articless, nil

}

func (h *Handler) GenerateQuiz(c *gin.Context) {
	var body CreateQuizDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("[quiz] GenerateQuiz bind error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[quiz] GenerateQuiz payload - categories=%v difficulty=%s questions=%d", body.CategoryIDs, body.Difficulty, body.NumberofQuestions)
	articles, err := h.GetArticles(c, body.CategoryIDs)
	if err != nil {
		log.Printf("[quiz] GenerateQuiz GetArticles error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[quiz] GenerateQuiz using %d articles", len(articles))
	quiz, err := services.GenerateQuiz(body.NumberofQuestions, body.CategoryIDs, h.cfg.GROQ_API_KEY, body.Difficulty, articles)
	if err != nil {
		log.Printf("[quiz] GenerateQuiz LLM error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[quiz] GenerateQuiz success - title=%s questions=%d", quiz.Title, len(quiz.Questions))
	c.JSON(http.StatusOK, quiz)
}
