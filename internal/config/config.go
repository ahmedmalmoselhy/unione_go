package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on environment variables")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dbURL := ""
	if dbHost != "" {
		dbURL = "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=UTC"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if smtpPort == 0 {
		smtpPort = 587
	}

	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        port,
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    smtpPort,
		SMTPUser:    os.Getenv("SMTP_USER"),
		SMTPPass:    os.Getenv("SMTP_PASS"),
		SMTPFrom:    os.Getenv("SMTP_FROM"),
	}
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("database configuration is missing (DB_HOST, DB_USER, etc.)")
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("invalid PORT: %v", err)
	}

	// Warn if SMTP is missing but don't fail as it might be optional for some environments
	if c.SMTPHost == "" {
		log.Println("Warning: SMTP_HOST is not set, email features will be disabled")
	}

	return nil
}
