package services

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
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
	ExportEmployees(facultyID *uint) ([]byte, error)
}
...
func (s *employeeService) ExportEmployees(facultyID *uint) ([]byte, error) {
	employees, err := s.ListEmployees(facultyID)
	if err != nil {
		return nil, err
	}
	return s.impExpSvc.ExportToExcel(employees, "Employees")
}

type employeeService struct {
	userRepo  repository.UserRepository
	impExpSvc ImportExportService
}

func NewEmployeeService(userRepo repository.UserRepository, impExpSvc ImportExportService) EmployeeService {
	return &employeeService{userRepo: userRepo, impExpSvc: impExpSvc}
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
	rows, err := s.impExpSvc.ReadExcel(file, "Sheet1")
	if err != nil {
		return 0, err
	}

	successCount := 0
	// Expected headers: National ID, Email, First Name, Last Name, Phone, Gender, Date of Birth (YYYY-MM-DD), Department ID
	for i, row := range rows {
		if i == 0 || len(row) < 4 {
			continue
		}

		nationalID := row[0]
		email := row[1]
		firstName := row[2]
		lastName := row[3]

		if nationalID == "" || email == "" || firstName == "" || lastName == "" {
			continue
		}

		user := &models.User{
			NationalID: nationalID,
			Email:      email,
			FirstName:  firstName,
			LastName:   lastName,
			Role:       models.RoleStudent,
			FacultyID:  &facultyID,
		}

		// Optional fields
		if len(row) > 4 {
			user.Phone = row[4]
		}
		if len(row) > 5 {
			user.Gender = row[5]
		}
		if len(row) > 6 && row[6] != "" {
			dob, err := time.Parse("2006-01-02", row[6])
			if err == nil {
				user.DateOfBirth = &dob
			}
		}
		if len(row) > 7 && row[7] != "" {
			deptID, err := strconv.ParseUint(row[7], 10, 32)
			if err == nil {
				uDeptID := uint(deptID)
				user.DepartmentID = &uDeptID
			}
		}

		// Default password is National ID
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(nationalID), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		if err := s.userRepo.CreateUser(user); err == nil {
			successCount++
		}
	}

	return successCount, nil
}
