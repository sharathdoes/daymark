package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	DBUrl                string
	GROQ_API_KEY         string
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
	APP_BASE_URL         string
	GOOGLE_CALLBACK_URL  string
	GITHUB_CALLBACK_URL  string
	FRONTEND_URL         string
	JWT_SECRET           string
	SESSION_SECRET       string
	SMTP_HOST            string
	SMTP_PORT            string
	SMTP_USERNAME        string
	SMTP_PASSWORD        string
	SMTP_FROM            string
}

func Load() *Config {
	_ = godotenv.Load()
	port := getEnv("PORT", "8080")
	appBaseURL := getEnv("APP_BASE_URL", "http://localhost:"+port)

	return &Config{
		Port:                 port,
		DBUrl:                os.Getenv("DATABASE_URL"),
		GROQ_API_KEY:         os.Getenv("GROQ_API_KEY"),
		GOOGLE_CLIENT_ID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GITHUB_CLIENT_ID:     os.Getenv("GITHUB_CLIENT_ID"),
		GITHUB_CLIENT_SECRET: os.Getenv("GITHUB_CLIENT_SECRET"),
		APP_BASE_URL:         appBaseURL,
		GOOGLE_CALLBACK_URL:  getEnv("GOOGLE_CALLBACK_URL", appBaseURL+"/user/google/callback"),
		GITHUB_CALLBACK_URL:  getEnv("GITHUB_CALLBACK_URL", appBaseURL+"/user/github/callback"),
		FRONTEND_URL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
		JWT_SECRET:           os.Getenv("JWT_SECRET"),
		SESSION_SECRET:       getEnv("SESSION_SECRET", "super-secret"),
		SMTP_HOST:            os.Getenv("SMTP_HOST"),
		SMTP_PORT:            getEnv("SMTP_PORT", "587"),
		SMTP_USERNAME:        os.Getenv("SMTP_USERNAME"),
		SMTP_PASSWORD:        os.Getenv("SMTP_PASSWORD"),
		SMTP_FROM:            os.Getenv("SMTP_FROM"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
