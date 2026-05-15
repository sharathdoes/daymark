package quiz

import (
	"context"
	"daymark/config"
	"daymark/internal/models"
	"daymark/internal/modules/articles"
	"daymark/internal/modules/feedSource"
	"daymark/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Scheduler is a minimal interface so the handler can trigger generation on demand.
type Scheduler interface {
	RunNow(ctx context.Context) error
}

type Handler struct {
	repo           *Repository
	FeedService    *feedSource.Service
	ArticleService *articles.Service
	cfg            config.Config
	Scheduler      Scheduler // may be nil if scheduler not wired
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

// GetDailyQuiz returns today's pre-generated Quiz of the Day.
// It is public (no auth required).
func (h *Handler) GetDailyQuiz(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	quiz, err := h.repo.GetDailyQuizByDate(c.Request.Context(), today)
	if err != nil {
		log.Printf("[quiz] GetDailyQuiz not found for date=%s err=%v", today, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No quiz of the day yet — check back later or after 6 AM IST.",
			"date":  today,
		})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

// TriggerDailyQuiz is a protected endpoint that immediately runs the daily quiz generation.
// Useful for testing and backfilling. Protected by JWT auth.
func (h *Handler) TriggerDailyQuiz(c *gin.Context) {
	if h.Scheduler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "scheduler not configured"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Minute)
	defer cancel()

	if err := h.Scheduler.RunNow(ctx); err != nil {
		log.Printf("[quiz] TriggerDailyQuiz error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the newly created (or existing) quiz
	today := time.Now().Format("2006-01-02")
	quiz, err := h.repo.GetDailyQuizByDate(c.Request.Context(), today)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "daily quiz generated successfully"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "daily quiz generated successfully",
		"quiz":    quiz,
	})
}
