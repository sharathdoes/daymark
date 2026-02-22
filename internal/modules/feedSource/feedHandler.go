package feedSource

import (

	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

func (h *Handler) CreateFeed(c *gin.Context) {
	var body CreateFeedDTO
	if err:=c.ShouldBindJSON(&body); err!=nil {
		c.JSON(401, gin.H{"error":"Body doesn't Match the Model"})
		return
	}
	if err:=h.s.CreateFeed(c, body.Name, body.URL, body.CategoryId); err!=nil {
		c.JSON(500, gin.H{"error":err.Error()})
	}
	c.JSON(201, gin.H{"message":"feed Created Successfully"})
}



func (h *Handler) GetFeedSourcesByCategory(c *gin.Context) {
	var body CategoriesDTO
	if err:=c.ShouldBindJSON(&body); err!=nil {
		c.JSON(401, gin.H{"error":"Body doesn't Match the Model"})
		return
	}
	feed, err := h.s.GetFeedSourcesByCategory(c, body.CategoryIds)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, feed)
}

