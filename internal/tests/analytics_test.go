package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func TestAnalyticsEndpoints(t *testing.T) {
	r, db, _ := setupIntegrationTest(t)

	// 1. Setup Admin User
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("adminpass"), bcrypt.DefaultCost)
	admin := models.User{
		Email:      "admin@unione.com",
		Password:   string(hashedPassword),
		Role:       models.RoleAdmin,
		FirstName:  "Super",
		LastName:   "Admin",
		NationalID: "ADMIN777",
	}
	db.Create(&admin)

	// 2. Setup some data for analytics
	uni := models.University{Name: "Uni"}
	db.Create(&uni)
	fac := models.Faculty{Name: "Fac", UniversityID: uni.ID}
	db.Create(&fac)
	dept := models.Department{Name: "Dept", FacultyID: fac.ID}
	db.Create(&dept)
	course := models.Course{Code: "C1", Name: "Course 1", DepartmentID: dept.ID}
	db.Create(&course)
	
	student := models.User{
		Email:      "student@unione.com",
		Role:       models.RoleStudent,
		FacultyID:  &fac.ID,
		DepartmentID: &dept.ID,
		NationalID: "STU123",
	}
	db.Create(&student)

	// 3. Login to get token
	loginPayload := map[string]interface{}{
		"email":    "admin@unione.com",
		"password": "adminpass",
	}
	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["data"].(map[string]interface{})["token"].(string)

	// 4. Test Dashboard Stats
	req, _ = http.NewRequest("GET", "/api/v1/admin/dashboard/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Dashboard stats failed: %s", w.Body.String())
	}

	var statsResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &statsResp)
	overview := statsResp["data"].(map[string]interface{})["overview"].(map[string]interface{})
	if overview["students"].(float64) != 1 {
		t.Errorf("Expected 1 student, got %f", overview["students"].(float64))
	}

	// 5. Test Enrollment Trends
	req, _ = http.NewRequest("GET", "/api/v1/admin/analytics/enrollment-trends", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Enrollment trends failed: %s", w.Body.String())
	}

	// 6. Test Professor Workload
	req, _ = http.NewRequest("GET", "/api/v1/admin/analytics/professor-workload", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Professor workload failed: %s", w.Body.String())
	}
}
