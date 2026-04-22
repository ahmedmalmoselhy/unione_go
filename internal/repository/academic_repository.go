package repository

import (
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
)

type AcademicRepository interface {
	// Academic Term
	CreateTerm(term *models.AcademicTerm) error
	GetAllTerms() ([]models.AcademicTerm, error)
	GetTermByID(id uint) (*models.AcademicTerm, error)
	UpdateTerm(term *models.AcademicTerm) error
	DeleteTerm(id uint) error

	// Course
	CreateCourse(course *models.Course) error
	GetCoursesByDepartment(deptID uint) ([]models.Course, error)
	GetCourseByID(id uint) (*models.Course, error)
	SetCoursePrerequisites(courseID uint, prerequisiteIDs []uint) error
	UpdateCourse(course *models.Course) error
	DeleteCourse(id uint) error

	// Section
	CreateSection(section *models.Section) error
	GetSectionsByCourse(courseID uint) ([]models.Section, error)
	GetSectionsByTerm(termID uint) ([]models.Section, error)
	GetSectionsByProfessor(profID uint) ([]models.Section, error)
	GetSectionByID(id uint) (*models.Section, error)
	UpdateSection(section *models.Section) error
	DeleteSection(id uint) error

	// Enrollment
	CreateEnrollment(enrollment *models.Enrollment) error
	GetEnrollmentsByStudent(studentID uint) ([]models.Enrollment, error)
	GetEnrollmentsBySection(sectionID uint) ([]models.Enrollment, error)
	GetEnrollment(studentID, sectionID uint) (*models.Enrollment, error)
	UpdateEnrollment(enrollment *models.Enrollment) error
	DeleteEnrollment(id uint) error

	// Attendance
	CreateAttendance(attendance *models.Attendance) error
	GetAttendanceBySectionAndDate(sectionID uint, date time.Time) ([]models.Attendance, error)
	GetAttendanceByStudentAndSection(studentID, sectionID uint) ([]models.Attendance, error)

	// Exam
	CreateExam(exam *models.Exam) error
	GetExamsBySection(sectionID uint) ([]models.Exam, error)
}

type academicRepository struct {
	db *gorm.DB
}

func NewAcademicRepository(db *gorm.DB) AcademicRepository {
	return &academicRepository{db: db}
}

// Term implementations
func (r *academicRepository) CreateTerm(term *models.AcademicTerm) error {
	return r.db.Create(term).Error
}

func (r *academicRepository) GetAllTerms() ([]models.AcademicTerm, error) {
	var terms []models.AcademicTerm
	err := r.db.Find(&terms).Error
	return terms, err
}

func (r *academicRepository) GetTermByID(id uint) (*models.AcademicTerm, error) {
	var term models.AcademicTerm
	err := r.db.First(&term, id).Error
	return &term, err
}

func (r *academicRepository) UpdateTerm(term *models.AcademicTerm) error {
	return r.db.Save(term).Error
}

func (r *academicRepository) DeleteTerm(id uint) error {
	return r.db.Delete(&models.AcademicTerm{}, id).Error
}

// Course implementations
func (r *academicRepository) CreateCourse(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *academicRepository) GetCoursesByDepartment(deptID uint) ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Preload("Prerequisites").Where("department_id = ?", deptID).Find(&courses).Error
	return courses, err
}

func (r *academicRepository) GetCourseByID(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Prerequisites").First(&course, id).Error
	return &course, err
}

func (r *academicRepository) SetCoursePrerequisites(courseID uint, prerequisiteIDs []uint) error {
	course, err := r.GetCourseByID(courseID)
	if err != nil {
		return err
	}

	if len(prerequisiteIDs) == 0 {
		return r.db.Model(course).Association("Prerequisites").Clear()
	}

	var prerequisites []models.Course
	if err := r.db.Where("id IN ?", prerequisiteIDs).Find(&prerequisites).Error; err != nil {
		return err
	}
	if len(prerequisites) != len(prerequisiteIDs) {
		return gorm.ErrRecordNotFound
	}

	return r.db.Model(course).Association("Prerequisites").Replace(&prerequisites)
}

func (r *academicRepository) UpdateCourse(course *models.Course) error {
	return r.db.Save(course).Error
}

func (r *academicRepository) DeleteCourse(id uint) error {
	return r.db.Delete(&models.Course{}, id).Error
}

// Section implementations
func (r *academicRepository) CreateSection(section *models.Section) error {
	return r.db.Create(section).Error
}

func (r *academicRepository) GetSectionsByCourse(courseID uint) ([]models.Section, error) {
	var sections []models.Section
	err := r.db.Preload("Course").Preload("Course.Prerequisites").Preload("AcademicTerm").Preload("Professor").Where("course_id = ?", courseID).Find(&sections).Error
	return sections, err
}

func (r *academicRepository) GetSectionsByTerm(termID uint) ([]models.Section, error) {
	var sections []models.Section
	err := r.db.Preload("Course").Preload("Course.Prerequisites").Preload("AcademicTerm").Preload("Professor").Where("academic_term_id = ?", termID).Find(&sections).Error
	return sections, err
}

func (r *academicRepository) GetSectionsByProfessor(profID uint) ([]models.Section, error) {
	var sections []models.Section
	err := r.db.Preload("Course").Preload("AcademicTerm").Where("professor_id = ?", profID).Find(&sections).Error
	return sections, err
}

func (r *academicRepository) GetSectionByID(id uint) (*models.Section, error) {
	var section models.Section
	err := r.db.Preload("Course").Preload("Course.Prerequisites").Preload("AcademicTerm").Preload("Professor").First(&section, id).Error
	return &section, err
}

func (r *academicRepository) UpdateSection(section *models.Section) error {
	return r.db.Save(section).Error
}

func (r *academicRepository) DeleteSection(id uint) error {
	return r.db.Delete(&models.Section{}, id).Error
}

// Enrollment implementations
func (r *academicRepository) CreateEnrollment(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *academicRepository) GetEnrollmentsByStudent(studentID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Preload("Section.Course").Where("student_id = ?", studentID).Find(&enrollments).Error
	return enrollments, err
}

func (r *academicRepository) GetEnrollmentsBySection(sectionID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	err := r.db.Preload("Student").Where("section_id = ?", sectionID).Find(&enrollments).Error
	return enrollments, err
}

func (r *academicRepository) GetEnrollment(studentID, sectionID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Where("student_id = ? AND section_id = ?", studentID, sectionID).First(&enrollment).Error
	return &enrollment, err
}

func (r *academicRepository) UpdateEnrollment(enrollment *models.Enrollment) error {
	return r.db.Save(enrollment).Error
}

func (r *academicRepository) DeleteEnrollment(id uint) error {
	return r.db.Delete(&models.Enrollment{}, id).Error
}

// Attendance implementations
func (r *academicRepository) CreateAttendance(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *academicRepository) GetAttendanceBySectionAndDate(sectionID uint, date time.Time) ([]models.Attendance, error) {
	var attendance []models.Attendance
	err := r.db.Where("section_id = ? AND date = ?", sectionID, date).Find(&attendance).Error
	return attendance, err
}

func (r *academicRepository) GetAttendanceByStudentAndSection(studentID, sectionID uint) ([]models.Attendance, error) {
	var attendance []models.Attendance
	err := r.db.Where("student_id = ? AND section_id = ?", studentID, sectionID).Find(&attendance).Error
	return attendance, err
}

// Exam implementations
func (r *academicRepository) CreateExam(exam *models.Exam) error {
	return r.db.Create(exam).Error
}

func (r *academicRepository) GetExamsBySection(sectionID uint) ([]models.Exam, error) {
	var exams []models.Exam
	err := r.db.Preload("Section.Course").Where("section_id = ?", sectionID).Find(&exams).Error
	return exams, err
}
