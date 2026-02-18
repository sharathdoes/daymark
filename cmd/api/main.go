package main

import (
	"daymark/cmd/workers"
	"log"
)

func main() {
	// cfg := configs.Load()
	// server:=server.New(cfg)
	quiz, err := workers.GenerateQuiz()
	if err != nil {
		log.Print("error: ", err)
		return
	}
	log.Print(quiz)
	// server.Run()
}
