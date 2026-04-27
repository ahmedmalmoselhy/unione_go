package services

import (
	"encoding/json"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

func (s *AuditService) Log(userID uint, action, entityType, entityID string, oldValues, newValues interface{}, ip, userAgent string) error {
	var oldJSON, newJSON string

	if oldValues != nil {
		b, _ := json.Marshal(oldValues)
		oldJSON = string(b)
	}

	if newValues != nil {
		b, _ := json.Marshal(newValues)
		newJSON = string(b)
	}

	log := models.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValues:  oldJSON,
		NewValues:  newJSON,
		IPAddress:  ip,
		UserAgent:  userAgent,
	}

	return s.db.Create(&log).Error
}

func (s *AuditService) GetLogs(limit, offset int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	err := s.db.Model(&models.AuditLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Preload("User").Order("created_at desc").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (s *AuditService) GetLogsByEntity(entityType, entityID string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := s.db.Preload("User").Where("entity_type = ? AND entity_id = ?", entityType, entityID).Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (s *AuditService) GetLogsByUser(userID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&logs).Error
	return logs, err
}
