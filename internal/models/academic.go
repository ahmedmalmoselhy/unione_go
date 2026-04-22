package models

import (
	"time"

	"gorm.io/gorm"
)

type AcademicTerm struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"` // e.g., "Fall 2026"
	StartDate time.Time      `json:"start_date"`
	EndDate   time.Time      `json:"end_date"`
	IsActive  bool           `gorm:"default:false" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Course struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Code          string         `gorm:"uniqueIndex;not null" json:"code"` // e.g., "CS101"
	Name          string         `gorm:"not null" json:"name"`
	Credits       int            `gorm:"not null" json:"credits"`
	Description   string         `json:"description"`
	DepartmentID  uint           `gorm:"not null" json:"department_id"`
	Department    *Department    `json:"department,omitempty"`
	Prerequisites []Course       `gorm:"many2many:course_prerequisites;joinForeignKey:CourseID;joinReferences:PrerequisiteID" json:"prerequisites,omitempty"`
	Sections      []Section      `json:"sections,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type CoursePrerequisite struct {
	CourseID       uint `gorm:"primaryKey"`
	PrerequisiteID uint `gorm:"primaryKey"`
}

type Section struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CourseID       uint           `gorm:"not null" json:"course_id"`
	Course         *Course        `json:"course,omitempty"`
	AcademicTermID uint           `gorm:"not null" json:"academic_term_id"`
	AcademicTerm   *AcademicTerm  `json:"academic_term,omitempty"`
	ProfessorID    *uint          `json:"professor_id"` // Nullable if not assigned yet
	Professor      *User          `json:"professor,omitempty"`
	Capacity       int            `gorm:"not null" json:"capacity"`
	Schedule       string         `json:"schedule"` // e.g., "Mon/Wed 10:00-11:30"
	Enrollments    []Enrollment   `json:"enrollments,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Enrollment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	StudentID uint           `gorm:"uniqueIndex:idx_student_section;not null" json:"student_id"`
	Student   *User          `json:"student,omitempty"`
	SectionID uint           `gorm:"uniqueIndex:idx_student_section;not null" json:"section_id"`
	Section   *Section       `json:"section,omitempty"`
	Grade     *float64       `json:"grade"`                                             // Nullable until graded
	Status    string         `gorm:"type:varchar(20);default:'enrolled'" json:"status"` // e.g., enrolled, dropped, completed
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Attendance struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SectionID uint           `gorm:"not null" json:"section_id"`
	Section   *Section       `json:"section,omitempty"`
	StudentID uint           `gorm:"not null" json:"student_id"`
	Student   *User          `json:"student,omitempty"`
	Date      time.Time      `gorm:"not null" json:"date"`
	IsPresent bool           `gorm:"default:true" json:"is_present"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Exam struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SectionID uint           `gorm:"not null" json:"section_id"`
	Section   *Section       `json:"section,omitempty"`
	Date      time.Time      `gorm:"not null" json:"date"`
	Location  string         `json:"location"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
