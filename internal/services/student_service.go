package services

import (
	"errors"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type StudentService interface {
	CreateStudent(email, password, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error)
	ListStudents(facultyID, departmentID *uint) ([]models.User, error)
	GetStudent(id uint) (*models.User, error)
	UpdateStudent(id uint, firstName, lastName string, facultyID uint, departmentID *uint) (*models.User, error)
	DeleteStudent(id uint) error
}

type studentService struct {
	userRepo repository.UserRepository
	orgRepo  repository.OrgRepository
}

func NewStudentService(userRepo repository.UserRepository, orgRepo repository.OrgRepository) StudentService {
	return &studentService{userRepo: userRepo, orgRepo: orgRepo}
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
