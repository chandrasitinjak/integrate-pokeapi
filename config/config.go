package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL   string
	Port      string
	DBUser    string
	DBPass    string
	DBHost    string
	DBPort    string
	DBName    string
	RedisAddr string
	RedisPass string
	RedisDB   int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	return &Config{
		BaseURL:   os.Getenv("BASE_URL"),
		Port:      os.Getenv("PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_NAME"),
		RedisAddr: os.Getenv("REDIS_ADDR"),
		RedisPass: os.Getenv("REDIS_PASS"),
		RedisDB:   redisDB,
	}, nil
}
