package main

import (
	"daymark/config"
	server "daymark/internal/router"

	docs "daymark/docs"
)

// 	"daymark/cmd/workers"
// 	"log"

func main() {
	cfg := config.Load()
	docs.SwaggerInfo.Title = "Daymark API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "API documentation for Daymark services"
	docs.SwaggerInfo.BasePath = "/"
	server := server.New(cfg)
	// quiz, err := workers.GenerateQuiz()
	// if err != nil {
	// 	log.Print("error: ", err)
	// 	return
	// }
	// log.Print(quiz)
	server.Run()
}
