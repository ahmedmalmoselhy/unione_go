package main

import (
	"log"
	"os"

	"github.com/ahmedmalmoselhy/unione_go/internal/db"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Initialize DB
	db.InitDB()

	// Auto-migrate schema
	log.Println("Migrating database schemas...")
	err := db.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	// Setup Router
	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting UniOne backend on :%s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
