package repository

import (
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
)

type AnnouncementRepository interface {
	Create(announcement *models.Announcement) error
	FindByFaculty(facultyID uint) ([]models.Announcement, error)
	FindBySection(sectionID uint) ([]models.Announcement, error)
	FindGlobal() ([]models.Announcement, error)
}

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) Create(announcement *models.Announcement) error {
	return r.db.Create(announcement).Error
}

func (r *announcementRepository) FindByFaculty(facultyID uint) ([]models.Announcement, error) {
	var announcements []models.Announcement
	err := r.db.Preload("Author").Where("type = ? AND faculty_id = ?", models.AnnouncementFaculty, facultyID).Find(&announcements).Error
	return announcements, err
}

func (r *announcementRepository) FindBySection(sectionID uint) ([]models.Announcement, error) {
	var announcements []models.Announcement
	err := r.db.Preload("Author").Where("type = ? AND section_id = ?", models.AnnouncementSection, sectionID).Find(&announcements).Error
	return announcements, err
}

func (r *announcementRepository) FindGlobal() ([]models.Announcement, error) {
	var announcements []models.Announcement
	err := r.db.Preload("Author").Where("type = ?", models.AnnouncementGlobal).Find(&announcements).Error
	return announcements, err
}
