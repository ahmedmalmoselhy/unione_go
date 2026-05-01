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

func TestStudentPortal(t *testing.T) {
	r, db, _ := setupIntegrationTest(t)

	// 1. Setup Student
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("studentpass"), bcrypt.DefaultCost)
	student := models.User{
		Email:      "student@portal.com",
		Password:   string(hashedPassword),
		Role:       models.RoleStudent,
		FirstName:  "Student",
		LastName:   "User",
		NationalID: "STU999",
	}
	db.Create(&student)

	// Login
	loginPayload := map[string]interface{}{
		"email":    "student@portal.com",
		"password": "studentpass",
	}
	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["data"].(map[string]interface{})["token"].(string)

	// 2. Test Profile
	req, _ = http.NewRequest("GET", "/api/v1/student/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Student profile failed: %s", w.Body.String())
	}

	// 3. Test Transcript (empty)
	req, _ = http.NewRequest("GET", "/api/v1/student/transcript", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Student transcript failed: %s", w.Body.String())
	}
}

func TestProfessorPortal(t *testing.T) {
	r, db, _ := setupIntegrationTest(t)

	// 1. Setup Professor
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("profpass"), bcrypt.DefaultCost)
	professor := models.User{
		Email:      "prof@portal.com",
		Password:   string(hashedPassword),
		Role:       models.RoleProfessor,
		FirstName:  "Prof",
		LastName:   "User",
		NationalID: "PROF999",
	}
	db.Create(&professor)

	// Login
	loginPayload := map[string]interface{}{
		"email":    "prof@portal.com",
		"password": "profpass",
	}
	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["data"].(map[string]interface{})["token"].(string)

	// 2. Test Sections (empty)
	req, _ = http.NewRequest("GET", "/api/v1/professor/sections", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Professor sections failed: %s", w.Body.String())
	}
}
