package user

import (
	"daymark/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type Handler struct {
	service     *Service
	frontendURL string
}

func NewHandler(s *Service, cfg *config.Config) *Handler {
	frontendURL := "http://localhost:3000"
	if cfg != nil && cfg.FRONTEND_URL != "" {
		frontendURL = cfg.FRONTEND_URL
	}

	return &Handler{service: s, frontendURL: frontendURL}
}

func (h *Handler) SignUp(c *gin.Context) {

	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.SignUp(body.Name, body.Email, body.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) SignIn(c *gin.Context) {

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.SignIn(body.Email, body.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GoogleLogin(c *gin.Context) {

	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *Handler) GoogleCallback(c *gin.Context) {

	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()

	u, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	_, err = h.service.OAuthLogin(
		u.Provider,
		u.UserID,
		u.Name,
		u.Email,
		u.AvatarURL,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, h.frontendURL+"/dashboard")
}

func (h *Handler) GithubLogin(c *gin.Context) {

	q := c.Request.URL.Query()
	q.Add("provider", "github")
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}
func (h *Handler) GithubCallback(c *gin.Context) {

	q := c.Request.URL.Query()
	q.Add("provider", "github")
	c.Request.URL.RawQuery = q.Encode()

	u, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.OAuthLogin(
		u.Provider,
		u.UserID,
		u.Name,
		u.Email,
		u.AvatarURL,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}