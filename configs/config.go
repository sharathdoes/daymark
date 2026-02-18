package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBUrl string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{Port: os.Getenv("PORT"), DBUrl: os.Getenv("DB_URL")}
}