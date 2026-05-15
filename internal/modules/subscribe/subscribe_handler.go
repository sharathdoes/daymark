package subscribe

import (
	"daymark/internal/models"
	"daymark/pkg/email"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo        *Repository
	emailSender *email.Sender
}

func NewHandler(repo *Repository, sender *email.Sender) *Handler {
	return &Handler{repo: repo, emailSender: sender}
}

type subscribeRequest struct {
	Email       string `json:"email" binding:"required,email"`
	CategoryIDs []int  `json:"category_ids"`
}

func (h *Handler) Subscribe(c *gin.Context) {
	var body subscribeRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parts := make([]string, len(body.CategoryIDs))
	for i, id := range body.CategoryIDs {
		parts[i] = strconv.Itoa(id)
	}

	sub := &models.Subscriber{
		Email:       body.Email,
		CategoryIDs: strings.Join(parts, ","),
	}

	if err := h.repo.Upsert(c.Request.Context(), sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to subscribe"})
		return
	}

	go h.emailSender.SendOTP(body.Email, "Welcome to Daymark", "You've subscribed to Daymark daily news digests. Stay informed!")

	c.JSON(http.StatusOK, gin.H{"message": "Subscribed"})
}
