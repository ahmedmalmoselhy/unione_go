package routes

import (
	"net/http"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/config"
	"github.com/ahmedmalmoselhy/unione_go/internal/handlers"
	"github.com/ahmedmalmoselhy/unione_go/internal/middlewares"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestLogger())

	// Rate Limiter: 10 req/s, burst of 20
	limiter := middlewares.NewIPRateLimiter(rate.Limit(10), 20)
	r.Use(middlewares.RateLimiterMiddleware(limiter))

	// Dependency Injection
	auditSvc := services.NewAuditService(db)
	webhookSvc := services.NewWebhookService(db)

	impExpSvc := services.NewImportExportService()
	transcriptSvc := services.NewTranscriptService()

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	orgRepo := repository.NewOrgRepository(db)
	orgService := services.NewOrgService(orgRepo)
	orgHandler := handlers.NewOrgHandler(orgService)

	employeeService := services.NewEmployeeService(userRepo, impExpSvc)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	templateHandler := handlers.NewTemplateHandler(impExpSvc)

	studentService := services.NewStudentService(userRepo, orgRepo, impExpSvc, auditSvc, webhookSvc)
	studentHandler := handlers.NewStudentHandler(studentService)

	professorService := services.NewProfessorService(userRepo, orgRepo, impExpSvc, auditSvc, webhookSvc)
	professorHandler := handlers.NewProfessorHandler(professorService)

	annRepo := repository.NewAnnouncementRepository(db)
	academicRepo := repository.NewAcademicRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	notifSvc := services.NewNotificationService(annRepo, userRepo, academicRepo, notifRepo, cfg)
	annHandler := handlers.NewAnnouncementHandler(notifSvc)
	notifHandler := handlers.NewNotificationHandler(notifSvc)

	academicService := services.NewAcademicService(academicRepo, userRepo, notifSvc, impExpSvc, auditSvc, webhookSvc)
	academicHandler := handlers.NewAcademicHandler(academicService)

	portalService := services.NewStudentPortalService(userRepo, academicRepo, academicService, transcriptSvc)
	portalHandler := handlers.NewStudentPortalHandler(portalService)

	profPortalService := services.NewProfessorPortalService(userRepo, academicRepo, annRepo, academicService, notifSvc)
	profPortalHandler := handlers.NewProfessorPortalHandler(profPortalService)

	govHandler := handlers.NewGovernanceHandler(auditSvc, webhookSvc)
	healthHandler := handlers.NewHealthHandler(db)

	api := r.Group("/api")

	// Health check endpoint
	api.GET("/health", healthHandler.HealthCheck)

	v1 := api.Group("/v1")
	{
		v1.Use(middlewares.AuditMiddleware(auditSvc))

		v1.GET("/health", healthHandler.HealthCheck)

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/me", middlewares.AuthMiddleware(cfg.JWTSecret), authHandler.GetMe)
		}

		// Organizations
		orgs := v1.Group("/organizations", middlewares.AuthMiddleware(cfg.JWTSecret))
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
			orgs.GET("/faculties/:faculty_id/employees/export", middlewares.RequireRole("admin"), employeeHandler.ExportEmployees)
		}

		// Academic Catalog & Enrollments
		academic := v1.Group("/academic", middlewares.AuthMiddleware(cfg.JWTSecret))
		{
			academic.GET("/terms", academicHandler.GetTerms)
			academic.POST("/terms", middlewares.RequireRole("admin", "employee"), academicHandler.CreateTerm)
			academic.GET("/terms/:id", academicHandler.GetTerm)
			academic.PUT("/terms/:id", middlewares.RequireRole("admin", "employee"), academicHandler.UpdateTerm)
			academic.DELETE("/terms/:id", middlewares.RequireRole("admin"), academicHandler.DeleteTerm)

			academic.GET("/courses", academicHandler.GetCourses)
			academic.POST("/courses", middlewares.RequireRole("admin", "employee"), academicHandler.CreateCourse)
			academic.GET("/courses/:id", academicHandler.GetCourse)
			academic.PUT("/courses/:id", middlewares.RequireRole("admin", "employee"), academicHandler.UpdateCourse)
			academic.DELETE("/courses/:id", middlewares.RequireRole("admin"), academicHandler.DeleteCourse)

			academic.GET("/sections", academicHandler.GetSections)
			academic.POST("/sections", middlewares.RequireRole("admin", "employee"), academicHandler.CreateSection)
			academic.GET("/sections/:id", academicHandler.GetSection)
			academic.PUT("/sections/:id", middlewares.RequireRole("admin", "employee"), academicHandler.UpdateSection)
			academic.DELETE("/sections/:id", middlewares.RequireRole("admin"), academicHandler.DeleteSection)

			academic.POST("/sections/:section_id/enrollments", middlewares.RequireRole("admin", "employee"), academicHandler.Enroll)
			academic.GET("/sections/:section_id/enrollments", middlewares.RequireRole("admin", "employee", "professor"), academicHandler.GetSectionEnrollments)
			academic.DELETE("/sections/:section_id/enrollments/:student_id", middlewares.RequireRole("admin", "employee"), academicHandler.DropSection)
			academic.GET("/students/:student_id/enrollments", middlewares.RequireRole("admin", "employee", "student"), academicHandler.GetStudentEnrollments)
			academic.PATCH("/sections/:section_id/students/:student_id/grade", middlewares.RequireRole("admin", "professor"), academicHandler.UpdateGrade)

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

			// Exports
			academic.GET("/sections/:section_id/enrollments/export", middlewares.RequireRole("admin", "employee", "professor"), academicHandler.ExportEnrollments)
			academic.GET("/sections/:section_id/grades/export", middlewares.RequireRole("admin", "professor"), academicHandler.ExportGrades)
		}

		// Announcements
		announcements := v1.Group("/announcements", middlewares.AuthMiddleware(cfg.JWTSecret))
		{
			announcements.GET("/", annHandler.ListAnnouncements)
			announcements.POST("/", middlewares.RequireRole("admin", "employee", "professor"), annHandler.CreateAnnouncement)
			announcements.POST("/:id/read", annHandler.MarkRead)
		}

		// Notifications
		notifications := v1.Group("/notifications", middlewares.AuthMiddleware(cfg.JWTSecret))
		{
			notifications.GET("/", notifHandler.ListNotifications)
			notifications.POST("/read-all", notifHandler.MarkAllRead)
			notifications.POST("/:id/read", notifHandler.MarkRead)
			notifications.DELETE("/:id", notifHandler.Delete)
		}

		admin := v1.Group("/admin", middlewares.AuthMiddleware(cfg.JWTSecret), middlewares.RequireRole("admin", "employee"))
		{
			admin.GET("/students", studentHandler.ListStudents)
			admin.GET("/students/export", studentHandler.ExportStudents)
			admin.POST("/students", studentHandler.CreateStudent)
			admin.GET("/students/:id", studentHandler.GetStudent)
			admin.PUT("/students/:id", studentHandler.UpdateStudent)
			admin.DELETE("/students/:id", studentHandler.DeleteStudent)
			admin.POST("/students/:id/transfer", studentHandler.TransferStudent)
			admin.GET("/students/:id/transfers", studentHandler.GetTransferHistory)

			admin.GET("/professors", professorHandler.ListProfessors)
			admin.GET("/professors/export", professorHandler.ExportProfessors)
			admin.POST("/professors", professorHandler.CreateProfessor)
			admin.GET("/professors/:id", professorHandler.GetProfessor)
			admin.PUT("/professors/:id", professorHandler.UpdateProfessor)
			admin.DELETE("/professors/:id", professorHandler.DeleteProfessor)

			admin.GET("/employees", employeeHandler.ListEmployees)
			admin.GET("/employees/export", employeeHandler.ExportEmployees)
			admin.POST("/employees", employeeHandler.AdminCreateEmployee)
			admin.GET("/employees/:id", employeeHandler.GetEmployee)
			admin.PUT("/employees/:id", employeeHandler.AdminUpdateEmployee)
			admin.DELETE("/employees/:id", employeeHandler.DeleteEmployee)

			admin.GET("/academic/terms", academicHandler.GetTerms)
			admin.POST("/academic/terms", academicHandler.CreateTerm)
			admin.GET("/academic/terms/:id", academicHandler.GetTerm)
			admin.PUT("/academic/terms/:id", academicHandler.UpdateTerm)
			admin.DELETE("/academic/terms/:id", academicHandler.DeleteTerm)

			admin.GET("/academic/courses", academicHandler.GetCourses)
			admin.POST("/academic/courses", academicHandler.CreateCourse)
			admin.GET("/academic/courses/:id", academicHandler.GetCourse)
			admin.PUT("/academic/courses/:id", academicHandler.UpdateCourse)
			admin.DELETE("/academic/courses/:id", academicHandler.DeleteCourse)

			admin.GET("/academic/sections", academicHandler.GetSections)
			admin.POST("/academic/sections", academicHandler.CreateSection)
			admin.GET("/academic/sections/:id", academicHandler.GetSection)
			admin.PUT("/academic/sections/:id", academicHandler.UpdateSection)
			admin.DELETE("/academic/sections/:id", academicHandler.DeleteSection)
			admin.GET("/academic/sections/:section_id/teaching-assistants", academicHandler.ListTeachingAssistants)
			admin.POST("/academic/sections/:section_id/teaching-assistants", academicHandler.AssignTeachingAssistant)
			admin.DELETE("/academic/sections/:section_id/teaching-assistants/:id", academicHandler.RemoveTeachingAssistant)
			admin.GET("/academic/sections/:section_id/exam-schedule", academicHandler.GetExamSchedule)
			admin.POST("/academic/sections/:section_id/exam-schedule", academicHandler.CreateExamSchedule)
			admin.PATCH("/academic/sections/:section_id/exam-schedule", academicHandler.UpdateExamSchedule)
			admin.POST("/academic/sections/:section_id/exam-schedule/publish", academicHandler.PublishExamSchedule)
			admin.GET("/academic/sections/:section_id/group-projects", academicHandler.ListGroupProjects)
			admin.POST("/academic/sections/:section_id/group-projects", academicHandler.CreateGroupProject)
			admin.PATCH("/academic/sections/:section_id/group-projects/:project_id", academicHandler.UpdateGroupProject)
			admin.DELETE("/academic/sections/:section_id/group-projects/:project_id", academicHandler.DeleteGroupProject)
			admin.POST("/academic/sections/:section_id/group-projects/:project_id/members", academicHandler.AddGroupProjectMember)
			admin.DELETE("/academic/sections/:section_id/group-projects/:project_id/members/:member_id", academicHandler.RemoveGroupProjectMember)
			admin.GET("/import-templates/students", templateHandler.DownloadStudentsImportTemplate)
			admin.GET("/import-templates/grades", templateHandler.DownloadGradesImportTemplate)

			admin.GET("/organizations/universities", orgHandler.GetUniversities)
			admin.POST("/organizations/universities", orgHandler.CreateUniversity)
			admin.PUT("/organizations/universities/:id", orgHandler.UpdateUniversity)
			admin.DELETE("/organizations/universities/:id", orgHandler.DeleteUniversity)

			admin.GET("/organizations/faculties", orgHandler.ListFaculties)
			admin.POST("/organizations/faculties", orgHandler.AdminCreateFaculty)
			admin.PUT("/organizations/faculties/:id", orgHandler.AdminUpdateFaculty)
			admin.DELETE("/organizations/faculties/:id", orgHandler.DeleteFaculty)

			admin.GET("/organizations/departments", orgHandler.ListDepartments)
			admin.POST("/organizations/departments", orgHandler.AdminCreateDepartment)
			admin.PUT("/organizations/departments/:id", orgHandler.AdminUpdateDepartment)
			admin.DELETE("/organizations/departments/:id", orgHandler.DeleteDepartment)

			// Governance
			gov := admin.Group("/governance", middlewares.RequireRole("admin"))
			{
				gov.GET("/audit-logs", govHandler.GetAuditLogs)
				gov.GET("/webhooks", govHandler.ListWebhooks)
				gov.POST("/webhooks", govHandler.CreateWebhook)
				gov.DELETE("/webhooks/:id", govHandler.DeleteWebhook)
				gov.GET("/webhooks/:id/deliveries", govHandler.GetWebhookDeliveries)
			}
		}

		// Student Portal
		student := v1.Group("/student", middlewares.AuthMiddleware(cfg.JWTSecret), middlewares.RequireRole("student"))
		{
			student.GET("/profile", portalHandler.GetProfile)
			student.GET("/enrollments", portalHandler.GetEnrollments)
			student.GET("/transcript", portalHandler.GetTranscript)
			student.GET("/transcript/export", portalHandler.ExportTranscriptPDF)
			student.GET("/academic-history", portalHandler.GetAcademicHistory)
			student.GET("/schedule", portalHandler.GetSchedule)
			student.GET("/schedule/export", portalHandler.ExportScheduleICS)
			student.GET("/attendance", portalHandler.GetAttendanceSummary)

			// Waitlist
			student.GET("/waitlist", portalHandler.GetWaitlist)
			student.POST("/waitlist", portalHandler.JoinWaitlist)
			student.DELETE("/waitlist/:id", portalHandler.LeaveWaitlist)

			// Ratings
			student.POST("/courses/:course_id/rate", portalHandler.RateCourse)
		}

		// Professor Portal
		professor := v1.Group("/professor", middlewares.AuthMiddleware(cfg.JWTSecret), middlewares.RequireRole("professor"))
		{
			professor.GET("/profile", profPortalHandler.GetProfile)
			professor.GET("/sections", profPortalHandler.GetSections)
			professor.GET("/sections/:id/roster", profPortalHandler.GetSectionRoster)
			professor.PATCH("/sections/:id/grade", profPortalHandler.UpdateGrade)
			professor.POST("/sections/:id/attendance", profPortalHandler.RecordAttendance)
			professor.POST("/sections/:id/announcements", profPortalHandler.CreateSectionAnnouncement)
			professor.GET("/schedule", profPortalHandler.GetSchedule)
		}
	}

	return r
}
