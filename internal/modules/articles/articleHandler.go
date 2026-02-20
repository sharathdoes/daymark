package articles

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	serv *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{serv: s}
}

func (h *Handler) CreateArticlesOfCategories(c *gin.Context) {
	var body CreateArticlesOfCategoriesDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(401, gin.H{"error": "Body doesn't Match the Model"})
		return
	}
	go h.serv.CreateArticlesOfCategories(c, body.Categories)
	c.JSON(http.StatusAccepted, gin.H{
    "message": "Article fetch started",
	})

}
