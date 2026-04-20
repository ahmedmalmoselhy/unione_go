package repository

import (
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindEmployeesByFaculty(facultyID uint) ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindEmployeesByFaculty(facultyID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ? AND faculty_id = ?", models.RoleEmployee, facultyID).Find(&users).Error
	return users, err
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
