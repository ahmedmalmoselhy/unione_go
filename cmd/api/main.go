package main

import (
	"log"

	"github.com/ahmedmalmoselhy/unione_go/internal/config"
	"github.com/ahmedmalmoselhy/unione_go/internal/db"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/routes"
)

// @title UniOne API
// @version 1.0
// @description This is the UniOne Go backend server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

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
			&models.StudentDepartmentHistory{},
			&models.AcademicTerm{},
			&models.Course{},
			&models.CoursePrerequisite{},
			&models.Section{},
			&models.SectionTeachingAssistant{},
			&models.Enrollment{},
			&models.Attendance{},
			&models.Announcement{},
			&models.Exam{},
			&models.GroupProject{},
			&models.GroupProjectMember{},
			&models.Waitlist{},
			&models.CourseRating{},
			&models.Notification{},
			&models.AnnouncementRead{},
			&models.AuditLog{},
			&models.Webhook{},
			&models.WebhookDelivery{},
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
