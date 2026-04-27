package services

import (
	"errors"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type StudentService interface {
	CreateStudent(email, password, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error)
	ListStudents(facultyID, departmentID *uint) ([]models.User, error)
	GetStudent(id uint) (*models.User, error)
	UpdateStudent(id uint, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error)
	TransferStudent(id uint, facultyID uint, departmentID *uint) (*models.User, error)
	GetTransferHistory(id uint) ([]models.StudentDepartmentHistory, error)
	DeleteStudent(id uint) error
	ExportStudents(facultyID, departmentID *uint) ([]byte, error)
}

type studentService struct {
	userRepo   repository.UserRepository
	orgRepo    repository.OrgRepository
	impExpSvc  ImportExportService
	auditSvc   *AuditService
	webhookSvc *WebhookService
}

func NewStudentService(userRepo repository.UserRepository, orgRepo repository.OrgRepository, impExpSvc ImportExportService, auditSvc *AuditService, webhookSvc *WebhookService) StudentService {
	return &studentService{
		userRepo:   userRepo,
		orgRepo:    orgRepo,
		impExpSvc:  impExpSvc,
		auditSvc:   auditSvc,
		webhookSvc: webhookSvc,
	}
}

func (s *studentService) validateScope(facultyID uint, departmentID *uint) error {
	if _, err := s.orgRepo.GetFacultyByID(facultyID); err != nil {
		return errors.New("faculty not found")
	}

	if departmentID == nil {
		return nil
	}

	department, err := s.orgRepo.GetDepartmentByID(*departmentID)
	if err != nil {
		return errors.New("department not found")
	}
	if department.FacultyID != facultyID {
		return errors.New("department does not belong to the provided faculty")
	}

	return nil
}

func (s *studentService) CreateStudent(email, password, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error) {
	if err := s.validateScope(facultyID, departmentID); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		Password:     string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		Role:         models.RoleStudent,
		FacultyID:    &facultyID,
		DepartmentID: departmentID,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// Trigger webhook
	s.webhookSvc.Trigger("student.created", map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})

	return user, nil
}

func (s *studentService) ListStudents(facultyID, departmentID *uint) ([]models.User, error) {
	role := models.RoleStudent
	return s.userRepo.FindUsers(repository.UserFilter{
		Role:         &role,
		FacultyID:    facultyID,
		DepartmentID: departmentID,
	})
}

func (s *studentService) GetStudent(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user.Role != models.RoleStudent {
		return nil, errors.New("user is not a student")
	}
	return user, nil
}

func (s *studentService) UpdateStudent(id uint, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error) {
	if err := s.validateScope(facultyID, departmentID); err != nil {
		return nil, err
	}

	user, err := s.GetStudent(id)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.FacultyID = &facultyID
	user.DepartmentID = departmentID

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *studentService) DeleteStudent(id uint) error {
	user, err := s.GetStudent(id)
	if err != nil {
		return err
	}

	return s.userRepo.DeleteUser(user.ID)
}

func (s *studentService) TransferStudent(id uint, facultyID uint, departmentID *uint) (*models.User, error) {
	if err := s.validateScope(facultyID, departmentID); err != nil {
		return nil, err
	}

	user, err := s.GetStudent(id)
	if err != nil {
		return nil, err
	}

	history := &models.StudentDepartmentHistory{
		StudentID:        user.ID,
		FromDepartmentID: user.DepartmentID,
		ToDepartmentID:   departmentID,
		FromFacultyID:    user.FacultyID,
		ToFacultyID:      &facultyID,
		TransferredAt:    time.Now().UTC(),
	}

	user.FacultyID = &facultyID
	user.DepartmentID = departmentID

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	if err := s.userRepo.CreateStudentDepartmentHistory(history); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *studentService) GetTransferHistory(id uint) ([]models.StudentDepartmentHistory, error) {
	if _, err := s.GetStudent(id); err != nil {
		return nil, err
	}

	return s.userRepo.GetStudentDepartmentHistory(id)
}

func (s *studentService) ExportStudents(facultyID, departmentID *uint) ([]byte, error) {
	students, err := s.ListStudents(facultyID, departmentID)
	if err != nil {
		return nil, err
	}
	return s.impExpSvc.ExportToExcel(students, "Students")
}
