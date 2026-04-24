package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      *User          `json:"user,omitempty"`
	Type      string         `gorm:"not null" json:"type"`
	Title     string         `gorm:"not null" json:"title"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Data      string         `gorm:"type:text" json:"data"` // JSON data
	ReadAt    *time.Time     `json:"read_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type AnnouncementRead struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	AnnouncementID uint      `gorm:"not null;index" json:"announcement_id"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	ReadAt         time.Time `gorm:"not null" json:"read_at"`
}
