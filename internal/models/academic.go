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
	ID                 uint                       `gorm:"primaryKey" json:"id"`
	CourseID           uint                       `gorm:"not null" json:"course_id"`
	Course             *Course                    `json:"course,omitempty"`
	AcademicTermID     uint                       `gorm:"not null" json:"academic_term_id"`
	AcademicTerm       *AcademicTerm              `json:"academic_term,omitempty"`
	ProfessorID        *uint                      `json:"professor_id"` // Nullable if not assigned yet
	Professor          *User                      `json:"professor,omitempty"`
	TeachingAssistants []SectionTeachingAssistant `json:"teaching_assistants,omitempty"`
	Capacity           int                        `gorm:"not null" json:"capacity"`
	Schedule           string                     `json:"schedule"` // e.g., "Mon/Wed 10:00-11:30"
	Enrollments        []Enrollment               `json:"enrollments,omitempty"`
	CreatedAt          time.Time                  `json:"created_at"`
	UpdatedAt          time.Time                  `json:"updated_at"`
	DeletedAt          gorm.DeletedAt             `gorm:"index" json:"-"`
}

type SectionTeachingAssistant struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	SectionID        uint           `gorm:"uniqueIndex:idx_section_professor_ta;not null" json:"section_id"`
	Section          *Section       `json:"section,omitempty"`
	ProfessorID      uint           `gorm:"uniqueIndex:idx_section_professor_ta;not null" json:"professor_id"`
	Professor        *User          `json:"professor,omitempty"`
	AssignedByUserID *uint          `json:"assigned_by_user_id"`
	AssignedBy       *User          `gorm:"foreignKey:AssignedByUserID" json:"assigned_by,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type Enrollment struct {
	ID        uint           `gorm:"primaryKey" json:"id" excel:"ID"`
	StudentID uint           `gorm:"uniqueIndex:idx_student_section;not null" json:"student_id" excel:"Student ID"`
	Student   *User          `json:"student,omitempty" excel:"-"`
	SectionID uint           `gorm:"uniqueIndex:idx_student_section;not null" json:"section_id" excel:"Section ID"`
	Section   *Section       `json:"section,omitempty" excel:"-"`
	Grade     *float64       `json:"grade" excel:"Grade"`                               // Nullable until graded
	Status    string         `gorm:"type:varchar(20);default:'enrolled'" json:"status" excel:"Status"` // e.g., enrolled, dropped, completed
	CreatedAt time.Time      `json:"created_at" excel:"Created At"`
	UpdatedAt time.Time      `json:"updated_at" excel:"Updated At"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" excel:"-"`
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
	ID          uint           `gorm:"primaryKey" json:"id"`
	SectionID   uint           `gorm:"not null" json:"section_id"`
	Section     *Section       `json:"section,omitempty"`
	Date        time.Time      `gorm:"not null" json:"date"`
	StartTime   string         `json:"start_time"`
	EndTime     string         `json:"end_time"`
	Location    string         `json:"location"`
	IsPublished bool           `gorm:"default:false" json:"is_published"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type GroupProject struct {
	ID              uint                 `gorm:"primaryKey" json:"id"`
	SectionID       uint                 `gorm:"not null;index" json:"section_id"`
	Section         *Section             `json:"section,omitempty"`
	Title           string               `gorm:"not null" json:"title"`
	Description     string               `json:"description"`
	DueAt           *time.Time           `json:"due_at"`
	MaxMembers      int                  `gorm:"default:5" json:"max_members"`
	IsActive        bool                 `gorm:"default:true" json:"is_active"`
	CreatedByUserID *uint                `json:"created_by_user_id"`
	CreatedBy       *User                `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
	Members         []GroupProjectMember `json:"members,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	DeletedAt       gorm.DeletedAt       `gorm:"index" json:"-"`
}

type GroupProjectMember struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	GroupProjectID uint           `gorm:"uniqueIndex:idx_group_project_member;not null" json:"group_project_id"`
	GroupProject   *GroupProject  `json:"group_project,omitempty"`
	StudentID      uint           `gorm:"uniqueIndex:idx_group_project_member;not null" json:"student_id"`
	Student        *User          `json:"student,omitempty"`
	JoinedAt       time.Time      `json:"joined_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Waitlist struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	StudentID uint           `gorm:"uniqueIndex:idx_student_section_waitlist;not null" json:"student_id"`
	Student   *User          `json:"student,omitempty"`
	SectionID uint           `gorm:"uniqueIndex:idx_student_section_waitlist;not null" json:"section_id"`
	Section   *Section       `json:"section,omitempty"`
	Position  int            `gorm:"not null" json:"position"`
	Status    string         `gorm:"type:varchar(20);default:'waiting'" json:"status"` // waiting, promoted, cancelled
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CourseRating struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	StudentID uint           `gorm:"not null" json:"student_id"`
	Student   *User          `json:"student,omitempty"`
	CourseID  uint           `gorm:"not null" json:"course_id"`
	Course    *Course        `json:"course,omitempty"`
	Rating    int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
