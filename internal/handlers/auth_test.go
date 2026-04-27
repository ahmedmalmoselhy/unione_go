package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedmalmoselhy/unione_go/internal/config"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() (*gin.Engine, services.AuthService) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})

	userRepo := repository.NewUserRepository(db)
	cfg := &config.Config{JWTSecret: "test_secret"}
	authService := services.NewAuthService(userRepo, cfg)
	
	return r, authService
}

func TestAuthHandler_Register(t *testing.T) {
	r, authSvc := setupTestRouter()
	handler := NewAuthHandler(authSvc)
	r.POST("/register", handler.Register)

	payload := map[string]interface{}{
		"email":      "handler@example.com",
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
		"role":       "student",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", resp.Code, resp.Body.String())
	}
}

func TestAuthHandler_Login(t *testing.T) {
	r, authSvc := setupTestRouter()
	handler := NewAuthHandler(authSvc)
	r.POST("/login", handler.Login)
	r.POST("/register", handler.Register)

	// Register first
	regPayload := map[string]interface{}{
		"email":      "login@example.com",
		"password":   "password123",
		"first_name": "Login",
		"last_name":  "User",
		"role":       "student",
	}
	regBody, _ := json.Marshal(regPayload)
	regReq, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(regBody))
	r.ServeHTTP(httptest.NewRecorder(), regReq)

	t.Run("Valid Login", func(t *testing.T) {
		loginPayload := map[string]interface{}{
			"email":    "login@example.com",
			"password": "password123",
		}
		loginBody, _ := json.Marshal(loginPayload)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.Code)
		}

		var data map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &data)
		if _, ok := data["data"].(map[string]interface{})["token"]; !ok {
			t.Error("Token missing in response")
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		loginPayload := map[string]interface{}{
			"email":    "login@example.com",
			"password": "wrongpassword",
		}
		loginBody, _ := json.Marshal(loginPayload)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.Code)
		}
	})
}
