package user

import (
	"daymark/config"
	"daymark/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	service := NewService(repo)
	h := NewHandler(service, cfg)

	g := r.Group("/auth")
	{
		g.POST("/signup", h.SignUp)
		g.POST("/signin", h.SignIn)
		g.POST("/verify-email", h.VerifyEmail)

		g.GET("/google", h.GoogleLogin)
		g.GET("/google/callback", h.GoogleCallback)

		g.GET("/github", h.GithubLogin)
		g.GET("/github/callback", h.GithubCallback)
	}

	// Authenticated user routes
	protected := g.Group("")
	protected.Use(utils.AuthMiddleware(cfg.JWT_SECRET))
	{
		protected.GET("/me", h.Me)
	}
}
