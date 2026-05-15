package articles

import (
	"context"
	"daymark/internal/models"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CategoryLookup interface {
	ListCategories(ctx context.Context) ([]models.Category, error)
}

type Handler struct {
	svc     *Service
	catLook CategoryLookup
}

func NewHandler(svc *Service, catLook CategoryLookup) *Handler {
	return &Handler{svc: svc, catLook: catLook}
}

type ArticleResponse struct {
	ID       uint   `json:"id"`
	Headline string `json:"headline"`
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Category string `json:"category"`
	ReadTime int    `json:"readTime"`
	URL      string `json:"url"`
}

// buildSummary returns content truncated cleanly at a word boundary.
func buildSummary(content string, maxChars int) string {
	if len(content) <= maxChars {
		return content
	}
	truncated := content[:maxChars]
	if idx := strings.LastIndexAny(truncated, " \t\n"); idx > 0 {
		truncated = truncated[:idx]
	}
	return strings.TrimRight(truncated, ".,;:") + "…"
}

func (h *Handler) GetTodayArticles(c *gin.Context) {
	arts, err := h.svc.GetTodayArticles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cats, _ := h.catLook.ListCategories(c.Request.Context())
	catMap := make(map[uint]string, len(cats))
	for _, cat := range cats {
		catMap[cat.ID] = cat.Name
	}

	resp := make([]ArticleResponse, 0, len(arts))
	for _, a := range arts {
		summary := buildSummary(a.Content, 300)
		wordCount := len(strings.Fields(a.Content))
		readTime := int(math.Ceil(float64(wordCount) / 200.0))
		if readTime < 1 {
			readTime = 1
		}
		resp = append(resp, ArticleResponse{
			ID:       a.ID,
			Headline: a.Title,
			Summary:  summary,
			Source:   a.Source,
			Category: catMap[a.CategoryID],
			ReadTime: readTime,
			URL:      a.Link,
		})
	}

	c.JSON(http.StatusOK, resp)
}
