package repository

import (
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AnnouncementRepository interface {
	Create(announcement *models.Announcement) error
	FindByFaculty(facultyID uint) ([]models.Announcement, error)
	FindBySection(sectionID uint) ([]models.Announcement, error)
	FindGlobal() ([]models.Announcement, error)
	FindRelevantForUser(facultyID *uint, sectionIDs []uint) ([]models.Announcement, error)
	MarkAsRead(announcementID, userID uint) error
	GetReadStatus(userID uint) (map[uint]time.Time, error)
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

func (r *announcementRepository) FindRelevantForUser(facultyID *uint, sectionIDs []uint) ([]models.Announcement, error) {
	var announcements []models.Announcement

	query := r.db.Preload("Author").Where("type = ?", models.AnnouncementGlobal)

	if facultyID != nil {
		query = query.Or("type = ? AND faculty_id = ?", models.AnnouncementFaculty, *facultyID)
	}

	if len(sectionIDs) > 0 {
		query = query.Or("type = ? AND section_id IN ?", models.AnnouncementSection, sectionIDs)
	}

	err := query.Order("created_at desc").Find(&announcements).Error
	return announcements, err
}

func (r *announcementRepository) MarkAsRead(announcementID, userID uint) error {
	read := models.AnnouncementRead{
		AnnouncementID: announcementID,
		UserID:         userID,
		ReadAt:         time.Now(),
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "announcement_id"}, {Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"read_at": time.Now()}),
	}).Create(&read).Error
}

func (r *announcementRepository) GetReadStatus(userID uint) (map[uint]time.Time, error) {
	var reads []models.AnnouncementRead
	err := r.db.Where("user_id = ?", userID).Find(&reads).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint]time.Time)
	for _, read := range reads {
		result[read.AnnouncementID] = read.ReadAt
	}
	return result, nil
}
