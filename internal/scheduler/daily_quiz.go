package scheduler

import (
	"context"
	"daymark/config"
	"daymark/internal/models"
	"daymark/internal/modules/articles"
	"daymark/internal/modules/category"
	"daymark/internal/modules/feedSource"
	"daymark/internal/modules/quiz"
	"daymark/internal/services"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// DailyQuizScheduler orchestrates the cron-based Quiz of the Day generation.
type DailyQuizScheduler struct {
	db          *gorm.DB
	cfg         *config.Config
	quizRepo    *quiz.Repository
	catRepo     *category.Repository
	feedRepo    *feedSource.Repository
	articleRepo *articles.Repository
	c           *cron.Cron
}

// NewDailyQuizScheduler wires up all required repositories.
func NewDailyQuizScheduler(db *gorm.DB, cfg *config.Config) *DailyQuizScheduler {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("[scheduler] WARNING: could not load Asia/Kolkata timezone, falling back to UTC: %v", err)
		loc = time.UTC
	}
	return &DailyQuizScheduler{
		db:          db,
		cfg:         cfg,
		quizRepo:    quiz.NewRepository(db),
		catRepo:     category.NewRepository(db),
		feedRepo:    feedSource.NewRepository(db),
		articleRepo: articles.NewRepository(db),
		c:           cron.New(cron.WithLocation(loc)),
	}
}

// Start registers the 6 AM IST cron job and launches the cron runner.
// It is non-blocking — the cron library manages its own goroutines.
func (s *DailyQuizScheduler) Start() {
	_, err := s.c.AddFunc("0 6 * * *", func() {
		log.Println("[scheduler] Cron fired: generating Quiz of the Day")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := s.RunNow(ctx); err != nil {
			log.Printf("[scheduler] Quiz of the Day generation FAILED: %v", err)
		}
	})
	if err != nil {
		log.Printf("[scheduler] CRITICAL: failed to register cron job: %v", err)
		return
	}
	s.c.Start()
	log.Println("[scheduler] Daily quiz cron registered (6:00 AM IST, every day)")
}

// Stop gracefully stops the cron scheduler.
func (s *DailyQuizScheduler) Stop() {
	s.c.Stop()
	log.Println("[scheduler] Daily quiz cron stopped")
}

// RunNow generates and stores the Quiz of the Day for today.
// It is idempotent — if today's daily quiz already exists it returns nil immediately.
func (s *DailyQuizScheduler) RunNow(ctx context.Context) error {
	today := time.Now().Format("2006-01-02")
	log.Printf("[scheduler] RunNow start date=%s", today)

	// --- Idempotency check ---
	if s.quizRepo.HasDailyQuizForDate(ctx, today) {
		log.Printf("[scheduler] RunNow skipping — daily quiz already exists for %s", today)
		return nil
	}

	// --- 1. Fetch all categories ---
	cats, err := s.catRepo.GetAll(ctx)
	if err != nil {
		return err
	}
	if len(cats) == 0 {
		log.Println("[scheduler] RunNow: no categories found, skipping generation")
		return nil
	}

	categoryIDs := make([]uint, 0, len(cats))
	for _, c := range cats {
		categoryIDs = append(categoryIDs, c.ID)
	}
	log.Printf("[scheduler] RunNow categories=%v", categoryIDs)

	// --- 2. Collect articles across all categories ---
	// First try cached articles for today; if none, sync from RSS feeds.
	articleSvc := articles.NewService(s.articleRepo)
	feedSvc := feedSource.NewService(s.feedRepo)

	existingArticles, err := articleSvc.GetTodayArticlesByCategory(ctx, categoryIDs)
	if err != nil {
		log.Printf("[scheduler] RunNow: GetTodayArticlesByCategory error: %v (continuing)", err)
	}

	var allArticles []models.Article
	if len(existingArticles) > 0 {
		log.Printf("[scheduler] RunNow: using %d cached articles", len(existingArticles))
		allArticles = existingArticles
	} else {
		// Sync category-by-category so a failing feed in one category doesn't block others.
		for _, cat := range cats {
			catID := []uint{cat.ID}
			feeds, err := feedSvc.GetFeedSourcesByCategory(ctx, catID)
			if err != nil {
				log.Printf("[scheduler] RunNow: GetFeedSourcesByCategory skip cat=%s err=%v", cat.Name, err)
				continue
			}
			if len(feeds) == 0 {
				log.Printf("[scheduler] RunNow: no feeds for category=%s, skipping", cat.Name)
				continue
			}
			synced, err := articleSvc.SyncFromFeeds(ctx, feeds, catID)
			if err != nil {
				log.Printf("[scheduler] RunNow: SyncFromFeeds skip cat=%s err=%v", cat.Name, err)
				continue
			}
			log.Printf("[scheduler] RunNow: synced %d articles for category=%s", len(synced), cat.Name)
			allArticles = append(allArticles, synced...)
		}
	}

	if len(allArticles) == 0 {
		log.Println("[scheduler] RunNow: no articles available, aborting quiz generation")
		return nil
	}
	log.Printf("[scheduler] RunNow: total articles for generation=%d", len(allArticles))

	// --- 3. Generate the quiz ---
	generatedQuiz, err := services.GenerateQuiz(10, categoryIDs, s.cfg.GROQ_API_KEY, "medium", allArticles)
	if err != nil {
		return err
	}

	// --- 4. Save the quiz ---
	if err := s.quizRepo.SaveQuiz(ctx, generatedQuiz); err != nil {
		return err
	}

	// --- 5. Record the daily quiz row ---
	dq := &models.DailyQuiz{
		QuizID: generatedQuiz.ID,
		Date:   today,
	}
	if err := s.quizRepo.SaveDailyQuiz(ctx, dq); err != nil {
		return err
	}

	log.Printf("[scheduler] RunNow SUCCESS: daily quiz created id=%d title=%q questions=%d",
		generatedQuiz.ID, generatedQuiz.Title, len(generatedQuiz.Questions))
	return nil
}
