package repository

import (
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
)

type OrgRepository interface {
	CreateUniversity(u *models.University) error
	GetUniversities() ([]models.University, error)
	GetUniversityByID(id uint) (*models.University, error)
	UpdateUniversity(u *models.University) error
	DeleteUniversity(id uint) error
	
	CreateFaculty(f *models.Faculty) error
	GetFacultiesByUniversity(uniID uint) ([]models.Faculty, error)
	GetFacultyByID(id uint) (*models.Faculty, error)
	UpdateFaculty(f *models.Faculty) error
	DeleteFaculty(id uint) error
	
	CreateDepartment(d *models.Department) error
	GetDepartmentsByFaculty(facultyID uint) ([]models.Department, error)
	GetDepartmentByID(id uint) (*models.Department, error)
	UpdateDepartment(d *models.Department) error
	DeleteDepartment(id uint) error
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

func (r *orgRepository) GetUniversityByID(id uint) (*models.University, error) {
	var uni models.University
	if err := r.db.First(&uni, id).Error; err != nil {
		return nil, err
	}
	return &uni, nil
}

func (r *orgRepository) UpdateUniversity(u *models.University) error {
	return r.db.Save(u).Error
}

func (r *orgRepository) DeleteUniversity(id uint) error {
	return r.db.Delete(&models.University{}, id).Error
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

func (r *orgRepository) UpdateFaculty(f *models.Faculty) error {
	return r.db.Save(f).Error
}

func (r *orgRepository) DeleteFaculty(id uint) error {
	return r.db.Delete(&models.Faculty{}, id).Error
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

func (r *orgRepository) GetDepartmentByID(id uint) (*models.Department, error) {
	var dept models.Department
	if err := r.db.First(&dept, id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *orgRepository) UpdateDepartment(d *models.Department) error {
	return r.db.Save(d).Error
}

func (r *orgRepository) DeleteDepartment(id uint) error {
	return r.db.Delete(&models.Department{}, id).Error
}
