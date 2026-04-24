package repository

import (
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByUser(userID uint) ([]models.Notification, error)
	MarkAsRead(id, userID uint) error
	MarkAllAsRead(userID uint) error
	Delete(id, userID uint) error
	FindByID(id uint) (*models.Notification, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByUser(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(id, userID uint) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read_at", &now).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", &now).Error
}

func (r *notificationRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{}).Error
}

func (r *notificationRepository) FindByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}
