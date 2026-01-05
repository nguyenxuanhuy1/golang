package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	CloudinaryURL string
}

var Config AppConfig

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	Config = AppConfig{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPass:        os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		CloudinaryURL: os.Getenv("CLOUDINARY_URL"),
	}
}
