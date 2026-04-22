package services

import (
	"errors"
	"fmt"
	"io"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeService interface {
	CreateEmployee(email, password, firstName, lastName string, facultyID uint) (*models.User, error)
	GetEmployeesByFaculty(facultyID uint) ([]models.User, error)
	ListEmployees(facultyID *uint) ([]models.User, error)
	GetEmployee(id uint) (*models.User, error)
	UpdateEmployee(id uint, firstName, lastName string) (*models.User, error)
	UpdateEmployeeDetails(id uint, firstName, lastName string, facultyID uint) (*models.User, error)
	DeleteEmployee(id uint) error
	ImportStudentsFromExcel(file io.Reader, facultyID uint) (int, error)
}

type employeeService struct {
	userRepo repository.UserRepository
}

func NewEmployeeService(userRepo repository.UserRepository) EmployeeService {
	return &employeeService{userRepo: userRepo}
}

func (s *employeeService) CreateEmployee(email, password, firstName, lastName string, facultyID uint) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		Role:      models.RoleEmployee,
		FacultyID: &facultyID,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *employeeService) GetEmployeesByFaculty(facultyID uint) ([]models.User, error) {
	return s.userRepo.FindEmployeesByFaculty(facultyID)
}

func (s *employeeService) ListEmployees(facultyID *uint) ([]models.User, error) {
	role := models.RoleEmployee
	return s.userRepo.FindUsers(repository.UserFilter{
		Role:      &role,
		FacultyID: facultyID,
	})
}

func (s *employeeService) GetEmployee(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if user.Role != models.RoleEmployee {
		return nil, errors.New("user is not an employee")
	}

	return user, nil
}

func (s *employeeService) UpdateEmployee(id uint, firstName, lastName string) (*models.User, error) {
	user, err := s.GetEmployee(id)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *employeeService) UpdateEmployeeDetails(id uint, firstName, lastName string, facultyID uint) (*models.User, error) {
	user, err := s.GetEmployee(id)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.FacultyID = &facultyID

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *employeeService) DeleteEmployee(id uint) error {
	if _, err := s.GetEmployee(id); err != nil {
		return err
	}

	return s.userRepo.DeleteUser(id)
}

func (s *employeeService) ImportStudentsFromExcel(file io.Reader, facultyID uint) (int, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return 0, fmt.Errorf("failed to open excel reader: %v", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return 0, fmt.Errorf("failed to get rows: %v", err)
	}

	successCount := 0
	// Assume first row is header: Email, Password, First Name, Last Name, Department ID (optional)
	for i, row := range rows {
		if i == 0 || len(row) < 4 {
			continue
		}

		email := row[0]
		password := row[1]
		firstName := row[2]
		lastName := row[3]

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &models.User{
			Email:     email,
			Password:  string(hashedPassword),
			FirstName: firstName,
			LastName:  lastName,
			Role:      models.RoleStudent,
			FacultyID: &facultyID,
		}

		if err := s.userRepo.CreateUser(user); err == nil {
			successCount++
		}
	}

	return successCount, nil
}
