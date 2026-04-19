package repository

import (
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
)

type OrgRepository interface {
	CreateUniversity(u *models.University) error
	GetUniversities() ([]models.University, error)
	
	CreateFaculty(f *models.Faculty) error
	GetFacultiesByUniversity(uniID uint) ([]models.Faculty, error)
	GetFacultyByID(id uint) (*models.Faculty, error)
	
	CreateDepartment(d *models.Department) error
	GetDepartmentsByFaculty(facultyID uint) ([]models.Department, error)
}

type orgRepository struct {
	db *gorm.DB
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{db: db}
}

func (r *orgRepository) CreateUniversity(u *models.University) error {
	return r.db.Create(u).Error
}

func (r *orgRepository) GetUniversities() ([]models.University, error) {
	var unis []models.University
	if err := r.db.Find(&unis).Error; err != nil {
		return nil, err
	}
	return unis, nil
}

func (r *orgRepository) CreateFaculty(f *models.Faculty) error {
	return r.db.Create(f).Error
}

func (r *orgRepository) GetFacultiesByUniversity(uniID uint) ([]models.Faculty, error) {
	var faculties []models.Faculty
	if err := r.db.Where("university_id = ?", uniID).Find(&faculties).Error; err != nil {
		return nil, err
	}
	return faculties, nil
}

func (r *orgRepository) GetFacultyByID(id uint) (*models.Faculty, error) {
	var faculty models.Faculty
	if err := r.db.First(&faculty, id).Error; err != nil {
		return nil, err
	}
	return &faculty, nil
}

func (r *orgRepository) CreateDepartment(d *models.Department) error {
	return r.db.Create(d).Error
}

func (r *orgRepository) GetDepartmentsByFaculty(facultyID uint) ([]models.Department, error) {
	var depts []models.Department
	if err := r.db.Where("faculty_id = ?", facultyID).Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}
