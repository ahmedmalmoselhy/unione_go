package models

import (
	"time"

	"gorm.io/gorm"
)

type AnnouncementType string

const (
	AnnouncementGlobal  AnnouncementType = "global"
	AnnouncementFaculty AnnouncementType = "faculty"
	AnnouncementSection AnnouncementType = "section"
)

type Announcement struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	Title     string           `gorm:"not null" json:"title"`
	Content   string           `gorm:"type:text;not null" json:"content"`
	Type      AnnouncementType `gorm:"type:varchar(20);not null" json:"type"`
	FacultyID *uint            `json:"faculty_id"`
	SectionID *uint            `json:"section_id"`
	AuthorID  uint             `gorm:"not null" json:"author_id"`
	Author    *User            `json:"author,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}
