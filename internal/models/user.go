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
	ID                  uint                       `gorm:"primaryKey" json:"id"`
	Email               string                     `gorm:"uniqueIndex;not null" json:"email"`
	Password            string                     `gorm:"not null" json:"-"`
	FirstName           string                     `json:"first_name"`
	LastName            string                     `json:"last_name"`
	Role                Role                       `gorm:"type:varchar(20);not null" json:"role"`
	NationalID          string                     `gorm:"uniqueIndex" json:"national_id"`
	Phone               string                     `json:"phone"`
	Gender              string                     `gorm:"type:varchar(10)" json:"gender"`
	DateOfBirth         *time.Time                 `json:"date_of_birth"`
	FacultyID           *uint                      `json:"faculty_id"`
	DepartmentID        *uint                      `json:"department_id"`
	Faculty             *Faculty                   `json:"-"`
	Department          *Department                `json:"-"`
	DepartmentTransfers []StudentDepartmentHistory `gorm:"foreignKey:StudentID" json:"department_transfers,omitempty"`
	CreatedAt           time.Time                  `json:"created_at"`
	UpdatedAt           time.Time                  `json:"updated_at"`
	DeletedAt           gorm.DeletedAt             `gorm:"index" json:"-"`
}

type StudentDepartmentHistory struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	StudentID        uint           `gorm:"not null;index" json:"student_id"`
	Student          *User          `json:"student,omitempty"`
	FromDepartmentID *uint          `json:"from_department_id"`
	ToDepartmentID   *uint          `json:"to_department_id"`
	FromFacultyID    *uint          `json:"from_faculty_id"`
	ToFacultyID      *uint          `json:"to_faculty_id"`
	TransferredAt    time.Time      `gorm:"not null" json:"transferred_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
