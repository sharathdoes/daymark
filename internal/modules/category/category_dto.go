package category


type createCategoryDTO struct {
	Name string `json:"name" binding:"required"`
}