package server

import (
	"daymark/config"
	"daymark/internal/modules/category"
	"daymark/internal/modules/feedSource"
	"daymark/internal/modules/quiz"
	"daymark/pkg/database"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
	Config *config.Config
}

func New(cfg *config.Config) *Server {
	r := gin.Default()
	r.Use(cors.Default())

	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// API routes
	feedSource.RegisterRoutes(r, db, cfg)
	category.RegisterRoutes(r, db, cfg)
	quiz.RegisterRoutes(r, db, cfg)

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

	return &Server{Engine: r, Config: cfg}
}

func (s *Server) Run() error {
	if s.Config.Port == "" {
		s.Config.Port = "8080"
	}
	return s.Engine.Run(":" + s.Config.Port)
}
