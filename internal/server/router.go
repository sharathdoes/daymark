package server

import (
	"daymark/configs"
	"daymark/internal/modules/feedSource"
	"daymark/pkg/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
	Config *configs.Config
}

func New(cfg *configs.Config) *Server {
	r := gin.Default()
	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Print("Error in Connecting with Database")
	}
	feedSource.RegisterRoutes(r, db, cfg)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return &Server{Engine: r, Config: cfg}
}

func (s *Server) Run() error {
	return s.Engine.Run(":" + s.Config.Port)
}
