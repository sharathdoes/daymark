package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBUrl string
	GROQ_API_KEY string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{Port: os.Getenv("PORT"), DBUrl: os.Getenv("DB_URL"), GROQ_API_KEY: os.Getenv("GROQ_API_KEY")}
}


func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}