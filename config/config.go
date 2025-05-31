package config

import (
	"fmt"
	"os"
	"log"

	"github.com/joho/godotenv"
)

type Config struct{
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	SSLMode string
}

func LoadConfig() *Config{
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error with loading db")
	}
	return &Config{
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("SSL_MODE"),
	}
}

func (c *Config) DataName() string{
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,c.DBPassword,c.DBHost,c.DBPort,c.DBName,c.SSLMode,
	)
}
