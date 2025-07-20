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
	// Загружаем переменные из файла .env (или config.env)
	err := godotenv.Load("./config.env")
	if err != nil {
		log.Println("Не удалось загрузить config.env, будут использованы значения по умолчанию")
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
	log.Printf("Переменная окружения %s не найдена, используется значение по умолчанию", key)
	return defaultVal
}
