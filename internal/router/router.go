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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Print("Error in Connecting with Database")
	}
	feedSource.RegisterRoutes(r, db, cfg)
	category.RegisterRoutes(r, db, cfg)
	quiz.RegisterRoutes(r, db, cfg)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return &Server{Engine: r, Config: cfg}
}

func (s *Server) Run() error {
	return s.Engine.Run(":" + s.Config.Port)
}
