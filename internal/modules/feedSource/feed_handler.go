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

// CreateFeed godoc
// @Summary      Create feed source
// @Description  Create a new feed source with one or more category IDs
// @Tags         feed
// @Accept       json
// @Produce      json
// @Param        payload  body      CreateFeedDTO  true  "Feed payload"
// @Success      201      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /feed/create [post]
func (h *Handler) CreateFeed(c *gin.Context) {
	var body CreateFeedDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(401, gin.H{"error": "Body doesn't Match the Model"})
		return
	}
	// basic validation: at least one categoryId must be provided
	if len(body.CategoryIds) == 0 {
		c.JSON(400, gin.H{"error": "at least one categoryId is required"})
		return
	}

	if err := h.s.CreateFeed(c, body.Name, body.URL, body.CategoryIds); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "feed Created Successfully"})
}

// GetFeedSourcesByCategory godoc
// @Summary      Get feeds by category IDs
// @Description  Get feed sources that belong to any of the given category IDs
// @Tags         feed
// @Accept       json
// @Produce      json
// @Param        payload  body      CategoriesDTO      true  "Category IDs"
// @Success      200      {array}   models.FeedSource
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /feed/ofCategories [get]
func (h *Handler) GetFeedSourcesByCategory(c *gin.Context) {
	var body CategoriesDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(401, gin.H{"error": "Body doesn't Match the Model"})
		return
	}
	feed, err := h.s.GetFeedSourcesByCategory(c, body.CategoryIds)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, feed)
}
