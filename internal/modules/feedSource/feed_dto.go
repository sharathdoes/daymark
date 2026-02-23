package feedSource


type CreateFeedDTO struct {
	Name     string         `json:"name"`
	URL      string         `json:"url"`
	CategoryId uint  `json:"categoryId"`
}

type CategoriesDTO struct {
	CategoryIds []uint `json:"categoryId"`
}

