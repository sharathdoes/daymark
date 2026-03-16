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

	// Save the quiz to the database so it can be shared via ID
	if err := h.repo.SaveQuiz(c.Request.Context(), quiz); err != nil {
		log.Printf("[quiz] GenerateQuiz failed to save to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save quiz"})
		return
	}
	log.Printf("[quiz] GenerateQuiz after-save quiz.ID=%d", quiz.ID)

	c.JSON(http.StatusOK, quiz)
}

// GetQuizByID fetches a generated quiz by its numeric ID
func (h *Handler) GetQuizByID(c *gin.Context) {
	quizID := c.Param("id")
	if quizID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter required"})
		return
	}

	quiz, err := h.repo.GetQuizByID(c.Request.Context(), quizID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

// SaveResult saves a completed quiz result for the authenticated user
func (h *Handler) SaveResult(c *gin.Context) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var body SaveQuizResultDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := &models.UserQuizResult{
		UserID:         userID,
		QuizID:         body.QuizID,
		Score:          body.Score,
		TotalQuestions: body.TotalQuestions,
		Difficulty:     body.Difficulty,
		Categories:     body.Categories,
	}

	err := h.repo.SaveUserResult(c.Request.Context(), result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save result: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetResults fetches past completed quiz results for the authenticated user
func (h *Handler) GetResults(c *gin.Context) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	results, err := h.repo.GetResultsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load results: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
