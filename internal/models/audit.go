package models

import (
	"time"
)

type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	User       *User     `json:"user,omitempty"`
	Action     string    `gorm:"type:varchar(100);not null" json:"action"`
	EntityType string    `gorm:"type:varchar(100);not null" json:"entity_type"`
	EntityID   string    `gorm:"type:varchar(100)" json:"entity_id"`
	OldValues  string    `gorm:"type:text" json:"old_values"`
	NewValues  string    `gorm:"type:text" json:"new_values"`
	IPAddress  string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  string    `gorm:"type:text" json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}
