package feedSource

type CreateFeedDTO struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	CategoryIds []uint `json:"categoryIds"`
}

type CategoriesDTO struct {
	CategoryIds []uint `json:"categoryId"`
}
