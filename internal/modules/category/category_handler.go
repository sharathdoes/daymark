package category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) *Handler {
	return &Handler{srv: srv}
}

// GetCategories godoc
// @Summary      List categories
// @Description  Get all available categories
// @Tags         category
// @Produce      json
// @Success      200  {array}   models.Category
// @Failure      500  {object}  map[string]string
// @Router       /category/ [get]
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.srv.ListCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Get a single category by its ID
// @Tags         category
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /category/{id} [get]
func (h *Handler) GetCategoryByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	category, err := h.srv.GetCategoryByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}


// CreateCategory godoc
// @Summary      Create category
// @Description  Create a new category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        payload  body      createCategoryDTO  true  "Category payload"
// @Success      201      {object}  models.Category
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /category/ [post]
// CreateCategory creates a new category from JSON body
func (h *Handler) CreateCategory(c *gin.Context) {
	var body createCategoryDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body doesn't Match the Model"})
		return
	}

	category, err := h.srv.CreateCategory(c, body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}
