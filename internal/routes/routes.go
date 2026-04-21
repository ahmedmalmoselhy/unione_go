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

	employeeService := services.NewEmployeeService(userRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	annRepo := repository.NewAnnouncementRepository(db)
	academicRepo := repository.NewAcademicRepository(db)
	
	notifSvc := services.NewNotificationService(annRepo, userRepo, academicRepo, cfg)
	annHandler := handlers.NewAnnouncementHandler(notifSvc)

	academicService := services.NewAcademicService(academicRepo, userRepo, notifSvc)
	academicHandler := handlers.NewAcademicHandler(academicService)

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
			// Universities
			orgs.GET("/universities", orgHandler.GetUniversities)
			orgs.POST("/universities", middlewares.RequireRole("admin"), orgHandler.CreateUniversity)
			orgs.PUT("/universities/:id", middlewares.RequireRole("admin"), orgHandler.UpdateUniversity)
			orgs.DELETE("/universities/:id", middlewares.RequireRole("admin"), orgHandler.DeleteUniversity)
			
			// Faculties
			orgs.GET("/universities/:university_id/faculties", orgHandler.GetFaculties)
			orgs.POST("/universities/:university_id/faculties", middlewares.RequireRole("admin"), orgHandler.CreateFaculty)
			orgs.PUT("/faculties/:id", middlewares.RequireRole("admin"), orgHandler.UpdateFaculty)
			orgs.DELETE("/faculties/:id", middlewares.RequireRole("admin"), orgHandler.DeleteFaculty)
			
			// Departments
			orgs.GET("/faculties/:faculty_id/departments", orgHandler.GetDepartments)
			orgs.POST("/faculties/:faculty_id/departments", middlewares.RequireRole("admin", "employee"), middlewares.RequireFacultyScope(), orgHandler.CreateDepartment)
			orgs.PUT("/departments/:id", middlewares.RequireRole("admin", "employee"), middlewares.RequireFacultyScope(), orgHandler.UpdateDepartment)
			orgs.DELETE("/departments/:id", middlewares.RequireRole("admin", "employee"), middlewares.RequireFacultyScope(), orgHandler.DeleteDepartment)

			// Employees (Staff)
			orgs.GET("/faculties/:faculty_id/employees", middlewares.RequireRole("admin", "employee"), middlewares.RequireFacultyScope(), employeeHandler.GetEmployees)
			orgs.POST("/faculties/:faculty_id/employees", middlewares.RequireRole("admin"), employeeHandler.CreateEmployee)
			orgs.PUT("/employees/:id", middlewares.RequireRole("admin"), employeeHandler.UpdateEmployee)
			orgs.DELETE("/employees/:id", middlewares.RequireRole("admin"), employeeHandler.DeleteEmployee)

			// Import Students
			orgs.POST("/faculties/:faculty_id/students/import", middlewares.RequireRole("admin", "employee"), middlewares.RequireFacultyScope(), employeeHandler.ImportStudents)
		}

		// Academic Catalog & Enrollments
		academic := v1.Group("/academic", middlewares.AuthMiddleware())
		{
			// ... (keep existing routes)
			// Attendance
			academic.POST("/sections/:section_id/attendance", middlewares.RequireRole("admin", "professor"), academicHandler.RecordAttendance)
			academic.GET("/sections/:section_id/attendance", middlewares.RequireRole("admin", "professor"), academicHandler.GetSectionAttendance)

			// GPA
			academic.GET("/students/:student_id/gpa", middlewares.RequireRole("admin", "employee", "student"), academicHandler.GetGPA)

			// Import Grades
			academic.POST("/sections/:section_id/grades/import", middlewares.RequireRole("admin", "professor"), academicHandler.ImportGrades)

			// Exams
			academic.POST("/sections/:section_id/exams", middlewares.RequireRole("admin", "employee", "professor"), academicHandler.CreateExam)
			academic.GET("/sections/:section_id/exams", academicHandler.GetExams)
		}

		// Announcements
		announcements := v1.Group("/announcements", middlewares.AuthMiddleware())
		{
			announcements.POST("/", middlewares.RequireRole("admin", "employee", "professor"), annHandler.CreateAnnouncement)
		}
	}

	return r
}
