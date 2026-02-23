package quiz

type CreateQuizDTO struct {
	Difficulty string `json:"diffifculty" binding:"required"`
	CategoryIDs []uint `json:"category_ids"`
	NumberofQuestions int `json:"number_of_questions"`
}