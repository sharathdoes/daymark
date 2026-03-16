package user

import (
	"daymark/config"
	"daymark/internal/models"
	"daymark/pkg/email"
	"daymark/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type Handler struct {
	service     *Service
	frontendURL string
	jwtSecret   string
	emailSender *email.Sender
}

func NewHandler(s *Service, cfg *config.Config) *Handler {
	frontendURL := "http://localhost:3000"
	if cfg != nil && cfg.FRONTEND_URL != "" {
		frontendURL = cfg.FRONTEND_URL
	}

	jwtSecret := ""
	if cfg != nil {
		jwtSecret = cfg.JWT_SECRET
	}

	// Direct email sender, similar to your NotificationProcessor example.
	emailSender := &email.Sender{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: "587",
		Username: "justabountyhunter935@gmail.com",
		Password: "cvnnegxwyvxoxdom",
		From:     "justabountyhunter935@gmail.com",
	}

	return &Handler{service: s, frontendURL: frontendURL, jwtSecret: jwtSecret, emailSender: emailSender}
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

	user, otp, err := h.service.SignUp(body.Name, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Temporary debug log so you can see the OTP in server logs during development.
	log.Printf("debug: OTP for %s is %s", user.Email, otp)

	// Send OTP via email (best-effort). In dev you can still log the code.
	if h.emailSender != nil {
		go h.emailSender.SendOTP(user.Email, "Your Daymark verification code", fmt.Sprintf("Your verification code is %s", otp))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "otp_sent",
		"email":   user.Email,
	})
}

// VerifyEmail confirms an email + OTP pair for manual signups, then issues a JWT.
func (h *Handler) VerifyEmail(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.VerifyEmail(body.Email, body.OTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
		return
	}

	token, err := utils.GenerateJWT(*user, h.jwtSecret, 30*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
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
		// Surface specific auth errors (including "email not verified")
		status := http.StatusUnauthorized
		if err.Error() == "email not verified" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if h.jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
		return
	}

	token, err := utils.GenerateJWT(*user, h.jwtSecret, 30*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
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

	// Issue JWT token and redirect back to frontend callback with token
	if h.jwtSecret == "" {
		c.JSON(500, gin.H{"error": "JWT secret not configured"})
		return
	}

	token, err := utils.GenerateJWT(*user, h.jwtSecret, 30*24*time.Hour)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", h.frontendURL, url.QueryEscape(token))
	c.Redirect(http.StatusFound, redirectURL)
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

	if h.jwtSecret == "" {
		c.JSON(500, gin.H{"error": "JWT secret not configured"})
		return
	}

	token, err := utils.GenerateJWT(*user, h.jwtSecret, 30*24*time.Hour)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", h.frontendURL, url.QueryEscape(token))
	c.Redirect(http.StatusFound, redirectURL)
}

func (h *Handler) Me(c *gin.Context) {
	// Prefer email from JWT claims when available
	emailVal, _ := c.Get("email")
	email, _ := emailVal.(string)

	var user *models.User
	var err error

	if email != "" {
		user, err = h.service.GetByEmail(email)
	}

	// Fallback to userID from JWT if email is missing or lookup failed
	if user == nil || err != nil {
		userIDVal, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, ok := userIDVal.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in token"})
			return
		}

		user, err = h.service.GetByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, user)
}
