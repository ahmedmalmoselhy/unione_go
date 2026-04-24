package services

import (
	"log"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/config"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"github.com/wneessen/go-mail"
)

type NotificationService interface {
	CreateAnnouncement(announcement *models.Announcement) error
	SendEmail(to []string, subject, body string)
	NotifyStudentsInSections(sectionIDs []uint, subject, body string)
	CreateNotification(userID uint, notifType, title, message string, data string) error
	GetNotifications(userID uint, unreadOnly bool) ([]models.Notification, int, error)
	MarkAsRead(id, userID uint) error
	MarkAllAsRead(userID uint) error
	DeleteNotification(id, userID uint) error
	ListAnnouncements(userID uint) ([]map[string]interface{}, error)
	MarkAnnouncementRead(announcementID, userID uint) error
}

type notificationService struct {
	annRepo   repository.AnnouncementRepository
	userRepo  repository.UserRepository
	acadRepo  repository.AcademicRepository
	notifRepo repository.NotificationRepository
	cfg       *config.Config
}

func NewNotificationService(
	annRepo repository.AnnouncementRepository,
	userRepo repository.UserRepository,
	acadRepo repository.AcademicRepository,
	notifRepo repository.NotificationRepository,
	cfg *config.Config,
) NotificationService {
	return &notificationService{
		annRepo:   annRepo,
		userRepo:  userRepo,
		acadRepo:  acadRepo,
		notifRepo: notifRepo,
		cfg:       cfg,
	}
}

func (s *notificationService) CreateAnnouncement(announcement *models.Announcement) error {
	if err := s.annRepo.Create(announcement); err != nil {
		return err
	}

	// For section-scoped announcements, notify students via email and in-app notification
	if announcement.Type == models.AnnouncementSection && announcement.SectionID != nil {
		go func() {
			enrollments, err := s.acadRepo.GetEnrollmentsBySection(*announcement.SectionID)
			if err != nil {
				log.Printf("Failed to get enrollments for section %d: %v", *announcement.SectionID, err)
				return
			}

			var emails []string
			for _, e := range enrollments {
				if e.StudentID != 0 {
					// Create in-app notification
					s.CreateNotification(e.StudentID, "section_announcement", announcement.Title, announcement.Content, "")

					// Collect emails for batch send
					// We need to fetch user to get email if not in enrollment
					// But e.Student might be preloaded
					if e.Student != nil {
						emails = append(emails, e.Student.Email)
					}
				}
			}

			if len(emails) > 0 {
				s.SendEmail(emails, announcement.Title, announcement.Content)
			}
		}()
	}

	return nil
}

func (s *notificationService) CreateNotification(userID uint, notifType, title, message string, data string) error {
	notification := &models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		Data:    data,
	}
	return s.notifRepo.Create(notification)
}

func (s *notificationService) GetNotifications(userID uint, unreadOnly bool) ([]models.Notification, int, error) {
	notifications, err := s.notifRepo.FindByUser(userID)
	if err != nil {
		return nil, 0, err
	}

	var result []models.Notification
	unreadCount := 0
	for _, n := range notifications {
		if n.ReadAt == nil {
			unreadCount++
		}
		if !unreadOnly || n.ReadAt == nil {
			result = append(result, n)
		}
	}

	return result, unreadCount, nil
}

func (s *notificationService) MarkAsRead(id, userID uint) error {
	return s.notifRepo.MarkAsRead(id, userID)
}

func (s *notificationService) MarkAllAsRead(userID uint) error {
	return s.notifRepo.MarkAllAsRead(userID)
}

func (s *notificationService) DeleteNotification(id, userID uint) error {
	return s.notifRepo.Delete(id, userID)
}

func (s *notificationService) ListAnnouncements(userID uint) ([]map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	var facultyID *uint
	if user.FacultyID != nil {
		facultyID = user.FacultyID
	}

	// For students, get section IDs
	var sectionIDs []uint
	if user.Role == models.RoleStudent {
		enrollments, err := s.acadRepo.GetEnrollmentsByStudent(userID)
		if err == nil {
			for _, e := range enrollments {
				sectionIDs = append(sectionIDs, e.SectionID)
			}
		}
	}

	announcements, err := s.annRepo.FindRelevantForUser(facultyID, sectionIDs)
	if err != nil {
		return nil, err
	}

	readStatus, err := s.annRepo.GetReadStatus(userID)
	if err != nil {
		readStatus = make(map[uint]time.Time)
	}

	var result []map[string]interface{}
	for _, ann := range announcements {
		_, isRead := readStatus[ann.ID]
		authorName := ""
		if ann.Author != nil {
			authorName = ann.Author.FirstName + " " + ann.Author.LastName
		}

		result = append(result, map[string]interface{}{
			"id":         ann.ID,
			"title":      ann.Title,
			"content":    ann.Content,
			"type":       ann.Type,
			"faculty_id": ann.FacultyID,
			"section_id": ann.SectionID,
			"author":     authorName,
			"is_read":    isRead,
			"created_at": ann.CreatedAt,
		})
	}

	return result, nil
}

func (s *notificationService) MarkAnnouncementRead(announcementID, userID uint) error {
	return s.annRepo.MarkAsRead(announcementID, userID)
}

func (s *notificationService) SendEmail(to []string, subject, body string) {
	if s.cfg.SMTPHost == "" {
		log.Println("SMTP configuration is missing, skipping email delivery")
		return
	}

	go func() {
		m := mail.NewMsg()
		if err := m.From(s.cfg.SMTPFrom); err != nil {
			log.Printf("Failed to set FROM: %v", err)
			return
		}
		if err := m.To(to...); err != nil {
			log.Printf("Failed to set TO: %v", err)
			return
		}
		m.Subject(subject)
		m.SetBodyString(mail.TypeTextPlain, body)

		c, err := mail.NewClient(s.cfg.SMTPHost, mail.WithPort(s.cfg.SMTPPort),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(s.cfg.SMTPUser), mail.WithPassword(s.cfg.SMTPPass))
		if err != nil {
			log.Printf("Failed to create mail client: %v", err)
			return
		}

		if err := c.DialAndSend(m); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}()
}

func (s *notificationService) NotifyStudentsInSections(sectionIDs []uint, subject, body string) {
	for _, sectionID := range sectionIDs {
		go func(sid uint) {
			enrollments, err := s.acadRepo.GetEnrollmentsBySection(sid)
			if err != nil {
				return
			}

			var emails []string
			for _, e := range enrollments {
				if e.StudentID != 0 {
					// Create in-app notification
					s.CreateNotification(e.StudentID, "academic_alert", subject, body, "")

					if e.Student != nil {
						emails = append(emails, e.Student.Email)
					}
				}
			}

			if len(emails) > 0 {
				s.SendEmail(emails, subject, body)
			}
		}(sectionID)
	}
}
