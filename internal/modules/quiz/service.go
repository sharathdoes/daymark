package quiz

import (
	"daymark/internal/models"
	"daymark/pkg/utils"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Service struct {
	repo       *Repository
	jobManager *JobManager
}

func NewService(repo *Repository, jm *JobManager) *Service {
	return &Service{
		repo:       repo,
		jobManager: jm,
	}
}

func (s *Service) StartQuizJob(feeds []models.FeedSource) string {
	jobID := uuid.New().String()
	progress := s.jobManager.CreateJob(jobID)

	go s.process(jobID, progress, feeds)

	return jobID
}

func (s *Service) process(jobID string, progress chan ProgressEvent, feeds []models.FeedSource) {
	defer close(progress)

	send := func(stage, msg string, data interface{}) {
		progress <- ProgressEvent{
			Stage:   stage,
			Message: msg,
			Data:    data,
		}
	}

	send("feeds", "Fetching feeds...", nil)

	articles, err := utils.FetchArticlesFromFeeds(feeds)
	if err != nil {
		send("error", "Failed to fetch articles", err.Error())
		return
	}

	send("articles", fmt.Sprintf("Parsed %d articles", len(articles)), nil)

	apiKey := os.Getenv("GROQ_API_KEY")
	var allQuestions []models.Question

	for i, article := range articles {
		if i >= 3 {
			break
		}

		send("llm", "Generating questions", article.Title)

		qs, err := utils.GenerateQuestions(article.Content, apiKey)
		if err != nil {
			send("llm_error", "Groq failed", article.Title)
			continue
		}

		allQuestions = append(allQuestions, qs...)
	}

	if len(allQuestions) == 0 {
		send("error", "No questions generated", nil)
		return
	}

	if err := s.repo.SaveQuiz(allQuestions); err != nil {
		send("error", "Failed to save quiz", err.Error())
		return
	}

	send("completed", "Quiz ready 🎯", len(allQuestions))
}