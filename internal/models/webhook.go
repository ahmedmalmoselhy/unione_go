package models

import (
	"time"

	"gorm.io/gorm"
)

type Webhook struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	URL       string         `gorm:"not null" json:"url"`
	Event     string         `gorm:"not null;index" json:"event"` // e.g. "student.created", "*"
	Secret    string         `gorm:"not null" json:"-"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	FailCount int            `gorm:"default:0" json:"fail_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type WebhookDelivery struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	WebhookID    uint      `gorm:"not null;index" json:"webhook_id"`
	Webhook      *Webhook  `json:"webhook,omitempty"`
	Event        string    `json:"event"`
	Payload      string    `gorm:"type:text" json:"payload"`
	StatusCode   int       `json:"status_code"`
	ResponseBody string    `gorm:"type:text" json:"response_body"`
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	Duration     int64     `json:"duration_ms"`
	CreatedAt    time.Time `json:"created_at"`
}
