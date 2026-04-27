package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedmalmoselhy/unione_go/internal/config"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupIntegrationTest(t *testing.T) (*gin.Engine, *gorm.DB, *config.Config) {
	gin.SetMode(gin.TestMode)
	
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrations handled by SetupRouter via the shared DB instance
	cfg := &config.Config{
		JWTSecret: "integration_test_secret",
		Port:      "8080",
	}

	r := routes.SetupRouter(db, cfg)
	err = db.AutoMigrate(
		&models.User{}, &models.University{}, &models.Faculty{},
		&models.Department{}, &models.AcademicTerm{}, &models.Course{},
		&models.Section{}, &models.Enrollment{}, &models.Attendance{},
		&models.Announcement{}, &models.AnnouncementRead{},
		&models.Exam{}, &models.GroupProject{}, &models.GroupProjectMember{},
		&models.Waitlist{}, &models.CourseRating{}, &models.Notification{},
		&models.AuditLog{}, &models.Webhook{}, &models.WebhookDelivery{},
		&models.SectionTeachingAssistant{}, &models.StudentDepartmentHistory{},
	)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return r, db, cfg
}

func TestFullEnrollmentFlow(t *testing.T) {
	r, db, _ := setupIntegrationTest(t)

	// 1. Setup Data (University, Faculty, Department, Term, Course, Section)
	uni := models.University{Name: "Test Uni"}
	db.Create(&uni)
	
	fac := models.Faculty{Name: "Test Fac", UniversityID: uni.ID}
	db.Create(&fac)
	
	dept := models.Department{Name: "Test Dept", FacultyID: fac.ID}
	db.Create(&dept)
	
	term := models.AcademicTerm{Name: "Fall 2026", IsActive: true}
	db.Create(&term)
	
	course := models.Course{Code: "CS101", Name: "Intro to CS", Credits: 3, DepartmentID: dept.ID}
	db.Create(&course)
	
	section := models.Section{CourseID: course.ID, AcademicTermID: term.ID, Capacity: 30, Schedule: "Mon 10am"}
	db.Create(&section)

	// 2. Register a Student
	regPayload := map[string]interface{}{
		"email":      "student@example.com",
		"password":   "password123",
		"first_name": "John",
		"last_name":  "Doe",
		"role":       "student",
	}
	body, _ := json.Marshal(regPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Registration failed: %s", w.Body.String())
	}

	// 3. Login to get token
	loginPayload := map[string]interface{}{
		"email":    "student@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)

	// 4. Register an Admin to perform enrollment (since currently RequireRole("admin", "employee") for enrollment)
	// Alternatively, we can just hack the DB for the admin user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)
	admin := models.User{
		Email:      "admin@example.com",
		Password:   string(hashedPassword),
		Role:       models.RoleAdmin,
		NationalID: "ADMIN123",
	}
	db.Create(&admin)
	
	// Login as admin
	loginPayload = map[string]interface{}{
		"email":    "admin@example.com",
		"password": "adminpassword",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	adminToken := loginResp["data"].(map[string]interface{})["token"].(string)

	// 5. Enroll student in section (Admin action)
	enrollUrl := fmt.Sprintf("/api/v1/academic/sections/%d/enrollments", section.ID)
	enrollPayload := map[string]interface{}{
		"student_id": 1, // Student ID from step 2
	}
	body, _ = json.Marshal(enrollPayload)
	req, _ = http.NewRequest("POST", enrollUrl, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Enrollment failed: %s", w.Body.String())
	}

	var enrollResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &enrollResp)
	
	// Data check
	data := enrollResp["data"].(map[string]interface{})
	if data["student_id"].(float64) != 1 {
		t.Errorf("Expected student_id 1, got %f", data["student_id"].(float64))
	}
}
