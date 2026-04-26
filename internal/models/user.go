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
	ID                  uint                       `gorm:"primaryKey" json:"id" excel:"ID"`
	Email               string                     `gorm:"uniqueIndex;not null" json:"email" excel:"Email"`
	Password            string                     `gorm:"not null" json:"-" excel:"-"`
	FirstName           string                     `json:"first_name" excel:"First Name"`
	LastName            string                     `json:"last_name" excel:"Last Name"`
	Role                Role                       `gorm:"type:varchar(20);not null" json:"role" excel:"Role"`
	NationalID          string                     `gorm:"uniqueIndex" json:"national_id" excel:"National ID"`
	Phone               string                     `json:"phone" excel:"Phone"`
	Gender              string                     `gorm:"type:varchar(10)" json:"gender" excel:"Gender"`
	DateOfBirth         *time.Time                 `json:"date_of_birth" excel:"Date of Birth"`
	FacultyID           *uint                      `json:"faculty_id" excel:"Faculty ID"`
	DepartmentID        *uint                      `json:"department_id" excel:"Department ID"`
	Faculty             *Faculty                   `json:"-" excel:"-"`
	Department          *Department                `json:"-" excel:"-"`
	DepartmentTransfers []StudentDepartmentHistory `gorm:"foreignKey:StudentID" json:"department_transfers,omitempty" excel:"-"`
	CreatedAt           time.Time                  `json:"created_at" excel:"Created At"`
	UpdatedAt           time.Time                  `json:"updated_at" excel:"Updated At"`
	DeletedAt           gorm.DeletedAt             `gorm:"index" json:"-" excel:"-"`
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
