package services

import (
	"errors"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
)

type ProfessorPortalService interface {
	GetProfile(professorID uint) (*models.User, error)
	GetSections(professorID uint) ([]models.Section, error)
	GetSectionRoster(professorID, sectionID uint) ([]models.Enrollment, error)
	UpdateGrade(professorID, sectionID, studentID uint, grade float64) (*models.Enrollment, error)
	RecordAttendance(professorID, sectionID uint, date time.Time, studentPresence map[uint]bool) error
	CreateSectionAnnouncement(professorID, sectionID uint, title, content string) (*models.Announcement, error)
	GetSchedule(professorID uint) (interface{}, error)
}

type professorPortalService struct {
	userRepo     repository.UserRepository
	academicRepo repository.AcademicRepository
	annRepo      repository.AnnouncementRepository
	academicSvc  AcademicService
	notifSvc     NotificationService
}

func NewProfessorPortalService(
	userRepo repository.UserRepository,
	academicRepo repository.AcademicRepository,
	annRepo repository.AnnouncementRepository,
	academicSvc AcademicService,
	notifSvc NotificationService,
) ProfessorPortalService {
	return &professorPortalService{
		userRepo:     userRepo,
		academicRepo: academicRepo,
		annRepo:      annRepo,
		academicSvc:  academicSvc,
		notifSvc:     notifSvc,
	}
}

func (s *professorPortalService) GetProfile(professorID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(professorID)
	if err != nil {
		return nil, err
	}
	if user.Role != models.RoleProfessor {
		return nil, errors.New("not a professor")
	}
	return user, nil
}

func (s *professorPortalService) GetSections(professorID uint) ([]models.Section, error) {
	return s.academicRepo.GetSectionsByProfessor(professorID)
}

func (s *professorPortalService) verifyOwnership(professorID, sectionID uint) error {
	section, err := s.academicRepo.GetSectionByID(sectionID)
	if err != nil {
		return err
	}
	if section.ProfessorID == nil || *section.ProfessorID != professorID {
		return errors.New("you do not teach this section")
	}
	return nil
}

func (s *professorPortalService) GetSectionRoster(professorID, sectionID uint) ([]models.Enrollment, error) {
	if err := s.verifyOwnership(professorID, sectionID); err != nil {
		return nil, err
	}
	return s.academicRepo.GetEnrollmentsBySection(sectionID)
}

func (s *professorPortalService) UpdateGrade(professorID, sectionID, studentID uint, grade float64) (*models.Enrollment, error) {
	if err := s.verifyOwnership(professorID, sectionID); err != nil {
		return nil, err
	}
	return s.academicSvc.UpdateGrade(studentID, sectionID, grade)
}

func (s *professorPortalService) RecordAttendance(professorID, sectionID uint, date time.Time, studentPresence map[uint]bool) error {
	if err := s.verifyOwnership(professorID, sectionID); err != nil {
		return err
	}
	return s.academicSvc.RecordAttendance(sectionID, date, studentPresence)
}

func (s *professorPortalService) CreateSectionAnnouncement(professorID, sectionID uint, title, content string) (*models.Announcement, error) {
	if err := s.verifyOwnership(professorID, sectionID); err != nil {
		return nil, err
	}

	announcement := &models.Announcement{
		Title:     title,
		Content:   content,
		AuthorID:  professorID,
		Type:      models.AnnouncementSection,
		SectionID: &sectionID,
	}

	if err := s.annRepo.Create(announcement); err != nil {
		return nil, err
	}

	// Notify students in this section
	s.notifSvc.NotifyStudentsInSections([]uint{sectionID}, title, content)

	return announcement, nil
}

func (s *professorPortalService) GetSchedule(professorID uint) (interface{}, error) {
	sections, err := s.academicRepo.GetSectionsByProfessor(professorID)
	if err != nil {
		return nil, err
	}

	var schedule []interface{}
	for _, sec := range sections {
		schedule = append(schedule, map[string]interface{}{
			"course_code": sec.Course.Code,
			"course_name": sec.Course.Name,
			"schedule":    sec.Schedule,
			"section_id":  sec.ID,
			"term":        sec.AcademicTerm.Name,
		})
	}

	return schedule, nil
}
