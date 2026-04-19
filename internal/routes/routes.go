package routes

import (
	"net/http"

	"github.com/ahmedmalmoselhy/unione_go/internal/config"
	"github.com/ahmedmalmoselhy/unione_go/internal/handlers"
	"github.com/ahmedmalmoselhy/unione_go/internal/middlewares"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	orgRepo := repository.NewOrgRepository(db)
	orgService := services.NewOrgService(orgRepo)
	orgHandler := handlers.NewOrgHandler(orgService)

	api := r.Group("/api")
	
	// Health check endpoint
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	v1 := api.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/me", middlewares.AuthMiddleware(), authHandler.GetMe)
		}

		// Organizations
		orgs := v1.Group("/organizations", middlewares.AuthMiddleware())
		{
			orgs.GET("/universities", orgHandler.GetUniversities)
			orgs.POST("/universities", middlewares.RequireRole("admin"), orgHandler.CreateUniversity)
			
			orgs.GET("/universities/:university_id/faculties", orgHandler.GetFaculties)
			orgs.POST("/universities/:university_id/faculties", middlewares.RequireRole("admin"), orgHandler.CreateFaculty) // Only admins make faculties
			
			orgs.GET("/faculties/:faculty_id/departments", orgHandler.GetDepartments)
			orgs.POST("/faculties/:faculty_id/departments", middlewares.RequireRole("admin", "employee"), middlewares.RequireFacultyScope(), orgHandler.CreateDepartment)
		}
	}

	return r
}
