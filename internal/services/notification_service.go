package services

import (
	"log"

	"github.com/ahmedmalmoselhy/unione_go/internal/config"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"github.com/wneessen/go-mail"
)

type NotificationService interface {
	CreateAnnouncement(announcement *models.Announcement) error
	SendEmail(to []string, subject, body string)
	NotifyStudentsInSections(sectionIDs []uint, subject, body string)
}

type notificationService struct {
	annRepo  repository.AnnouncementRepository
	userRepo repository.UserRepository
	acadRepo repository.AcademicRepository
	cfg      *config.Config
}

func NewNotificationService(
	annRepo repository.AnnouncementRepository,
	userRepo repository.UserRepository,
	acadRepo repository.AcademicRepository,
	cfg *config.Config,
) NotificationService {
	return &notificationService{
		annRepo:  annRepo,
		userRepo: userRepo,
		acadRepo: acadRepo,
		cfg:      cfg,
	}
}

func (s *notificationService) CreateAnnouncement(announcement *models.Announcement) error {
	if err := s.annRepo.Create(announcement); err != nil {
		return err
	}

	// For section-scoped announcements, notify students via email
	if announcement.Type == models.AnnouncementSection && announcement.SectionID != nil {
		go func() {
			enrollments, err := s.acadRepo.GetEnrollmentsBySection(*announcement.SectionID)
			if err != nil {
				log.Printf("Failed to get enrollments for section %d: %v", *announcement.SectionID, err)
				return
			}

			var emails []string
			for _, e := range enrollments {
				if e.Student != nil {
					emails = append(emails, e.Student.Email)
				}
			}

			if len(emails) > 0 {
				s.SendEmail(emails, announcement.Title, announcement.Content)
			}
		}()
	}

	return nil
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
				if e.Student != nil {
					emails = append(emails, e.Student.Email)
				}
			}

			if len(emails) > 0 {
				s.SendEmail(emails, subject, body)
			}
		}(sectionID)
	}
}
