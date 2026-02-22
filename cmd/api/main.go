package main

import (
	"daymark/config"
	"daymark/internal/router"
)

// 	"daymark/cmd/workers"
// 	"log"

func main() {
	cfg := config.Load()
	server:=server.New(cfg)
	// quiz, err := workers.GenerateQuiz()
	// if err != nil {
	// 	log.Print("error: ", err)
	// 	return
	// }
	// log.Print(quiz)
	server.Run()
}
