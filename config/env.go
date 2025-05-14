package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file")
	}
	log.Println(".env file loaded!")
}

type AuthConfig struct{
	Secret string
}

type DataBaseConfig struct {
	Url string
}

type LogConfig struct {
	Level  int
	Format string
}

func NewAuthConfig() *AuthConfig{
	return &AuthConfig{
		Secret: getString("SECRET", ""),
	}
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:  getInt("LOG_LEVEL", 0),
		Format: getString("LOG_FORMAT", "json"),
	}
}

func NewDBConfig() *DataBaseConfig {
	return &DataBaseConfig{
		Url: getString("DATABASE_URL", ""),
	}
}

func getString(key, defValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defValue
	}
	return value
}

func getInt(key string, defValue int) int {
	value := os.Getenv(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		return defValue
	}
	if value == "" {
		return defValue
	}
	return i
}

func getBool(key string, defValue bool) bool {
	value := os.Getenv(key)
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defValue
	}
	return b
}
