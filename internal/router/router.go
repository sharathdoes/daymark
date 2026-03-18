package server

import (
	"daymark/config"
	"daymark/internal/modules/category"
	"daymark/internal/modules/feedSource"
	"daymark/internal/modules/quiz"
	"daymark/internal/modules/user"
	"daymark/internal/scheduler"
	"daymark/pkg/database"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

type Server struct {
	Engine *gin.Engine
	Config *config.Config
}

func New(cfg *config.Config) *Server {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://daymark-eight.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	user.SetupProviders(cfg)

	store := sessions.NewCookieStore([]byte(cfg.SESSION_SECRET))
	gothic.Store = store

	// Start the daily quiz scheduler (6 AM IST every day)
	dailySched := scheduler.NewDailyQuizScheduler(db, cfg)
	dailySched.Start()

	// API routes
	feedSource.RegisterRoutes(r, db, cfg)
	category.RegisterRoutes(r, db, cfg)
	quiz.RegisterRoutes(r, db, cfg, dailySched)
	user.RegisterRoutes(r, db, cfg)

	r.GET("/debug-rss", func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			c.JSON(400, gin.H{"error": "url query param required"})
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		c.JSON(200, gin.H{
			"url":        url,
			"status":     resp.Status,
			"statusCode": resp.StatusCode,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return &Server{Engine: r, Config: cfg}
}

func (s *Server) Run() error {
	if s.Config.Port == "" {
		s.Config.Port = "8080"
	}
	return s.Engine.Run(":" + s.Config.Port)
}

