package services

import (
	"errors"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
)

type AcademicService interface {
	// Terms
	CreateTerm(name string, start, end time.Time) (*models.AcademicTerm, error)
	GetAllTerms() ([]models.AcademicTerm, error)
	GetTerm(id uint) (*models.AcademicTerm, error)
	UpdateTerm(id uint, name string, start, end time.Time, isActive bool) (*models.AcademicTerm, error)
	DeleteTerm(id uint) error

	// Courses
	CreateCourse(code, name string, credits int, desc string, deptID uint) (*models.Course, error)
	GetCoursesByDept(deptID uint) ([]models.Course, error)
	GetCourse(id uint) (*models.Course, error)
	UpdateCourse(id uint, code, name string, credits int, desc string) (*models.Course, error)
	DeleteCourse(id uint) error

	// Sections
	CreateSection(courseID, termID uint, capacity int, schedule string) (*models.Section, error)
	GetSectionsByCourse(courseID uint) ([]models.Section, error)
	GetSectionsByTerm(termID uint) ([]models.Section, error)
	GetSectionsByProfessor(profID uint) ([]models.Section, error)
	GetSection(id uint) (*models.Section, error)
	UpdateSection(id uint, capacity int, schedule string, professorID *uint) (*models.Section, error)
	DeleteSection(id uint) error

	// Enrollments
	EnrollStudent(studentID, sectionID uint) (*models.Enrollment, error)
	GetStudentEnrollments(studentID uint) ([]models.Enrollment, error)
	GetSectionEnrollments(sectionID uint) ([]models.Enrollment, error)
	UpdateGrade(studentID, sectionID uint, grade float64) (*models.Enrollment, error)
	DropSection(studentID, sectionID uint) error

	// Attendance
	RecordAttendance(sectionID uint, date time.Time, studentPresence map[uint]bool) error
	GetSectionAttendance(sectionID uint, date time.Time) ([]models.Attendance, error)
	GetStudentAttendance(studentID, sectionID uint) ([]models.Attendance, error)

	// GPA Calculation
	CalculateGPA(studentID uint) (float64, error)
}

type academicService struct {
	repo     repository.AcademicRepository
	userRepo repository.UserRepository
}

func NewAcademicService(repo repository.AcademicRepository, userRepo repository.UserRepository) AcademicService {
	return &academicService{repo: repo, userRepo: userRepo}
}

// Term logic
func (s *academicService) CreateTerm(name string, start, end time.Time) (*models.AcademicTerm, error) {
	term := &models.AcademicTerm{Name: name, StartDate: start, EndDate: end}
	if err := s.repo.CreateTerm(term); err != nil {
		return nil, err
	}
	return term, nil
}

func (s *academicService) GetAllTerms() ([]models.AcademicTerm, error) {
	return s.repo.GetAllTerms()
}

func (s *academicService) GetTerm(id uint) (*models.AcademicTerm, error) {
	return s.repo.GetTermByID(id)
}

func (s *academicService) UpdateTerm(id uint, name string, start, end time.Time, isActive bool) (*models.AcademicTerm, error) {
	term, err := s.repo.GetTermByID(id)
	if err != nil {
		return nil, err
	}
	term.Name = name
	term.StartDate = start
	term.EndDate = end
	term.IsActive = isActive
	if err := s.repo.UpdateTerm(term); err != nil {
		return nil, err
	}
	return term, nil
}

func (s *academicService) DeleteTerm(id uint) error {
	return s.repo.DeleteTerm(id)
}

// Course logic
func (s *academicService) CreateCourse(code, name string, credits int, desc string, deptID uint) (*models.Course, error) {
	course := &models.Course{Code: code, Name: name, Credits: credits, Description: desc, DepartmentID: deptID}
	if err := s.repo.CreateCourse(course); err != nil {
		return nil, err
	}
	return course, nil
}

func (s *academicService) GetCoursesByDept(deptID uint) ([]models.Course, error) {
	return s.repo.GetCoursesByDepartment(deptID)
}

func (s *academicService) GetCourse(id uint) (*models.Course, error) {
	return s.repo.GetCourseByID(id)
}

func (s *academicService) UpdateCourse(id uint, code, name string, credits int, desc string) (*models.Course, error) {
	course, err := s.repo.GetCourseByID(id)
	if err != nil {
		return nil, err
	}
	course.Code = code
	course.Name = name
	course.Credits = credits
	course.Description = desc
	if err := s.repo.UpdateCourse(course); err != nil {
		return nil, err
	}
	return course, nil
}

func (s *academicService) DeleteCourse(id uint) error {
	return s.repo.DeleteCourse(id)
}

// Section logic
func (s *academicService) CreateSection(courseID, termID uint, capacity int, schedule string) (*models.Section, error) {
	section := &models.Section{CourseID: courseID, AcademicTermID: termID, Capacity: capacity, Schedule: schedule}
	if err := s.repo.CreateSection(section); err != nil {
		return nil, err
	}
	return section, nil
}

func (s *academicService) GetSectionsByCourse(courseID uint) ([]models.Section, error) {
	return s.repo.GetSectionsByCourse(courseID)
}

func (s *academicService) GetSectionsByTerm(termID uint) ([]models.Section, error) {
	return s.repo.GetSectionsByTerm(termID)
}

func (s *academicService) GetSectionsByProfessor(profID uint) ([]models.Section, error) {
	return s.repo.GetSectionsByProfessor(profID)
}

func (s *academicService) GetSection(id uint) (*models.Section, error) {
	return s.repo.GetSectionByID(id)
}

func (s *academicService) UpdateSection(id uint, capacity int, schedule string, professorID *uint) (*models.Section, error) {
	section, err := s.repo.GetSectionByID(id)
	if err != nil {
		return nil, err
	}
	section.Capacity = capacity
	section.Schedule = schedule
	section.ProfessorID = professorID
	if err := s.repo.UpdateSection(section); err != nil {
		return nil, err
	}
	return section, nil
}

func (s *academicService) DeleteSection(id uint) error {
	return s.repo.DeleteSection(id)
}

// Enrollment logic
func (s *academicService) EnrollStudent(studentID, sectionID uint) (*models.Enrollment, error) {
	// 1. Check if student exists and has student role
	student, err := s.userRepo.FindByID(studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}
	if student.Role != models.RoleStudent {
		return nil, errors.New("user is not a student")
	}

	// 2. Check if section exists and has capacity
	section, err := s.repo.GetSectionByID(sectionID)
	if err != nil {
		return nil, errors.New("section not found")
	}

	var currentCount int64
	enrollments, _ := s.repo.GetEnrollmentsBySection(sectionID)
	currentCount = int64(len(enrollments))

	if currentCount >= int64(section.Capacity) {
		return nil, errors.New("section is at full capacity")
	}

	// 3. Create enrollment
	enrollment := &models.Enrollment{StudentID: studentID, SectionID: sectionID, Status: "enrolled"}
	if err := s.repo.CreateEnrollment(enrollment); err != nil {
		return nil, err
	}
	return enrollment, nil
}

func (s *academicService) GetStudentEnrollments(studentID uint) ([]models.Enrollment, error) {
	return s.repo.GetEnrollmentsByStudent(studentID)
}

func (s *academicService) GetSectionEnrollments(sectionID uint) ([]models.Enrollment, error) {
	return s.repo.GetEnrollmentsBySection(sectionID)
}

func (s *academicService) UpdateGrade(studentID, sectionID uint, grade float64) (*models.Enrollment, error) {
	enrollment, err := s.repo.GetEnrollment(studentID, sectionID)
	if err != nil {
		return nil, err
	}
	enrollment.Grade = &grade
	if err := s.repo.UpdateEnrollment(enrollment); err != nil {
		return nil, err
	}
	return enrollment, nil
}

func (s *academicService) DropSection(studentID, sectionID uint) error {
	enrollment, err := s.repo.GetEnrollment(studentID, sectionID)
	if err != nil {
		return err
	}
	return s.repo.DeleteEnrollment(enrollment.ID)
}

// Attendance logic
func (s *academicService) RecordAttendance(sectionID uint, date time.Time, studentPresence map[uint]bool) error {
	for studentID, isPresent := range studentPresence {
		attendance := &models.Attendance{
			SectionID: sectionID,
			StudentID: studentID,
			Date:      date,
			IsPresent: isPresent,
		}
		if err := s.repo.CreateAttendance(attendance); err != nil {
			return err
		}
	}
	return nil
}

func (s *academicService) GetSectionAttendance(sectionID uint, date time.Time) ([]models.Attendance, error) {
	return s.repo.GetAttendanceBySectionAndDate(sectionID, date)
}

func (s *academicService) GetStudentAttendance(studentID, sectionID uint) ([]models.Attendance, error) {
	return s.repo.GetAttendanceByStudentAndSection(studentID, sectionID)
}

// GPA Calculation logic
func (s *academicService) CalculateGPA(studentID uint) (float64, error) {
	enrollments, err := s.repo.GetEnrollmentsByStudent(studentID)
	if err != nil {
		return 0, err
	}

	var totalGradePoints float64
	var totalCredits int

	for _, e := range enrollments {
		if e.Grade != nil && e.Section != nil && e.Section.Course != nil {
			totalGradePoints += (*e.Grade) * float64(e.Section.Course.Credits)
			totalCredits += e.Section.Course.Credits
		}
	}

	if totalCredits == 0 {
		return 0, nil
	}

	return totalGradePoints / float64(totalCredits), nil
}
