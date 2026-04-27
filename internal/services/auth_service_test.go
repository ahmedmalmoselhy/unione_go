package services

import (
	"testing"

	"github.com/ahmedmalmoselhy/unione_go/internal/config"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestRegisterUser(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	cfg := &config.Config{JWTSecret: "test_secret"}
	svc := NewAuthService(repo, cfg)

	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
		Role:     models.RoleStudent,
	}

	err := svc.RegisterUser(user)
	if err != nil {
		t.Errorf("RegisterUser failed: %v", err)
	}

	// Verify user is in DB
	found, err := repo.FindByEmail("test@example.com")
	if err != nil {
		t.Errorf("Could not find user after registration: %v", err)
	}

	if found.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, found.Email)
	}
}

func TestAuthenticateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	cfg := &config.Config{JWTSecret: "test_secret"}
	svc := NewAuthService(repo, cfg)

	email := "auth@example.com"
	password := "secret123"
	user := &models.User{
		Email:    email,
		Password: password,
		Role:     models.RoleStudent,
	}

	_ = svc.RegisterUser(user)

	t.Run("Valid Credentials", func(t *testing.T) {
		token, err := svc.AuthenticateUser(email, password)
		if err != nil {
			t.Errorf("AuthenticateUser failed with valid credentials: %v", err)
		}
		if token == "" {
			t.Error("Expected token, got empty string")
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		_, err := svc.AuthenticateUser(email, "wrongpassword")
		if err == nil {
			t.Error("Expected error with invalid password, got nil")
		}
	})

	t.Run("Non-existent User", func(t *testing.T) {
		_, err := svc.AuthenticateUser("nonexistent@example.com", password)
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
	})
}
