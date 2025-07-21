package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath    string
	Port      string
	JWTSecret string
}

var AppConfig Config

func Load() {
	err := godotenv.Load("./config.env")
	if err != nil {
		log.Println("Failed to load config.env, default values will be used")
	}

	AppConfig = Config{
		DBPath:    getEnv("DB_PATH", "./store.db"),
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "default_secret_key"),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	log.Printf("Environment variable %s not found, using default value", key)
	return defaultVal
}
