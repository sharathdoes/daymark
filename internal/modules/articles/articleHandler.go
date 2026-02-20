package articles

import (
	"log"
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

	func(cat []string){
		if err := h.serv.CreateArticlesOfCategories(cat); err != nil {
			log.Println("background article creation failed:", err)
		}
	}(body.Categories)

	c.JSON(http.StatusAccepted, gin.H{
    "message": "Article fetch started",
	})

}
