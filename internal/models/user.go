package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleEmployee  Role = "employee"
	RoleProfessor Role = "professor"
	RoleStudent   Role = "student"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      Role           `gorm:"type:varchar(20);not null" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
