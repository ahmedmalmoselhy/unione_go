package services

import (
	"errors"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
)

type OrgService interface {
	CreateUniversity(name string) (*models.University, error)
	GetUniversities() ([]models.University, error)

	CreateFaculty(name string, uniID uint) (*models.Faculty, error)
	GetFaculties(uniID uint) ([]models.Faculty, error)

	CreateDepartment(name string, facultyID uint) (*models.Department, error)
	GetDepartments(facultyID uint) ([]models.Department, error)
}

type orgService struct {
	repo repository.OrgRepository
}

func NewOrgService(repo repository.OrgRepository) OrgService {
	return &orgService{repo: repo}
}

func (s *orgService) CreateUniversity(name string) (*models.University, error) {
	if name == "" {
		return nil, errors.New("university name cannot be empty")
	}
	uni := &models.University{Name: name}
	if err := s.repo.CreateUniversity(uni); err != nil {
		return nil, err
	}
	return uni, nil
}

func (s *orgService) GetUniversities() ([]models.University, error) {
	return s.repo.GetUniversities()
}

func (s *orgService) CreateFaculty(name string, uniID uint) (*models.Faculty, error) {
	if name == "" {
		return nil, errors.New("faculty name cannot be empty")
	}
	
	// Validate university exists indirectly or trust DB foreign key constraints.
	// For purity, you might want to fetch University first, but keeping it simple:
	faculty := &models.Faculty{Name: name, UniversityID: uniID}
	if err := s.repo.CreateFaculty(faculty); err != nil {
		return nil, err
	}
	return faculty, nil
}

func (s *orgService) GetFaculties(uniID uint) ([]models.Faculty, error) {
	return s.repo.GetFacultiesByUniversity(uniID)
}

func (s *orgService) CreateDepartment(name string, facultyID uint) (*models.Department, error) {
	if name == "" {
		return nil, errors.New("department name cannot be empty")
	}
	
	// Validate parent exists to provide a friendly error
	_, err := s.repo.GetFacultyByID(facultyID)
	if err != nil {
		return nil, errors.New("parent faculty not found")
	}

	dept := &models.Department{Name: name, FacultyID: facultyID}
	if err := s.repo.CreateDepartment(dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *orgService) GetDepartments(facultyID uint) ([]models.Department, error) {
	return s.repo.GetDepartmentsByFaculty(facultyID)
}
