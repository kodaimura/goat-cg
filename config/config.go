package config

import (
	"os"
	"log"
	"fmt"

	"github.com/joho/godotenv"
)


type Config struct {
	AppHost string
	AppPort string

	DBName string
	DBHost string
	DBPort string
	DBUser string
	DBPass string

	JwtSecretKey string
	LogLevel string
}

var cf Config


func init() {
	err := godotenv.Load(fmt.Sprintf("config/env/%s.env", os.Getenv("ENV")))

	if err != nil {
		log.Panic(err)
	}

	cf.AppHost = os.Getenv("APP_HOST")
	cf.AppPort = os.Getenv("APP_PORT")

	cf.DBName = os.Getenv("DB_NAME")
	cf.DBHost = os.Getenv("DB_HOST")
	cf.DBPort = os.Getenv("DB_PORT")
	cf.DBUser = os.Getenv("DB_USER")
	cf.DBPass = os.Getenv("DB_PASSWORD")

	cf.JwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	cf.LogLevel = os.Getenv("LOG_LEVEL")
}


func GetConfig() *Config{
	return &cf
}