package quiz

// Input payload for saving a quiz result
type SaveQuizResultDTO struct {
	QuizID         uint   `json:"quiz_id"`
	Score          int    `json:"score"`
	TotalQuestions int    `json:"total_questions"`
	Difficulty     string `json:"difficulty"`
	Categories     string `json:"categories"`
}
