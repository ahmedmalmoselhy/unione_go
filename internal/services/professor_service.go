package services

import (
	"errors"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type ProfessorService interface {
	CreateProfessor(email, password, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error)
	ListProfessors(facultyID, departmentID *uint) ([]models.User, error)
	GetProfessor(id uint) (*models.User, error)
	UpdateProfessor(id uint, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error)
	DeleteProfessor(id uint) error
	ExportProfessors(facultyID, departmentID *uint) ([]byte, error)
}

type professorService struct {
	userRepo  repository.UserRepository
	orgRepo   repository.OrgRepository
	impExpSvc ImportExportService
}

func NewProfessorService(userRepo repository.UserRepository, orgRepo repository.OrgRepository, impExpSvc ImportExportService) ProfessorService {
	return &professorService{userRepo: userRepo, orgRepo: orgRepo, impExpSvc: impExpSvc}
}

func (s *professorService) validateScope(facultyID uint, departmentID *uint) error {
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

func (s *professorService) CreateProfessor(email, password, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error) {
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
		Role:         models.RoleProfessor,
		FacultyID:    &facultyID,
		DepartmentID: departmentID,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *professorService) ListProfessors(facultyID, departmentID *uint) ([]models.User, error) {
	role := models.RoleProfessor
	return s.userRepo.FindUsers(repository.UserFilter{
		Role:         &role,
		FacultyID:    facultyID,
		DepartmentID: departmentID,
	})
}

func (s *professorService) GetProfessor(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user.Role != models.RoleProfessor {
		return nil, errors.New("user is not a professor")
	}
	return user, nil
}

func (s *professorService) UpdateProfessor(id uint, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error) {
	if err := s.validateScope(facultyID, departmentID); err != nil {
		return nil, err
	}

	user, err := s.GetProfessor(id)
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

func (s *professorService) DeleteProfessor(id uint) error {
	user, err := s.GetProfessor(id)
	if err != nil {
		return err
	}

	return s.userRepo.DeleteUser(user.ID)
}

func (s *professorService) ExportProfessors(facultyID, departmentID *uint) ([]byte, error) {
	professors, err := s.ListProfessors(facultyID, departmentID)
	if err != nil {
		return nil, err
	}
	return s.impExpSvc.ExportToExcel(professors, "Professors")
}
