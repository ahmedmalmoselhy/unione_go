package main

import (
	"log"

	"github.com/ahmedmalmoselhy/unione_go/internal/config"
	"github.com/ahmedmalmoselhy/unione_go/internal/db"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize DB
	db.InitDB(cfg.DatabaseURL)

	// Auto-migrate schema
	log.Println("Migrating database schemas...")
	if db.DB != nil {
		err := db.DB.AutoMigrate(
			&models.University{},
			&models.Faculty{},
			&models.Department{},
			&models.User{},
		)
		if err != nil {
			log.Fatalf("Failed to auto-migrate database: %v", err)
		}
	}

	// Setup Router
	r := routes.SetupRouter(db.DB, cfg)

	log.Printf("Starting UniOne backend on :%s...\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
