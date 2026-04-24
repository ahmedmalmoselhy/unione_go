package services

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type AcademicService interface {
	// Terms
	CreateTerm(name string, start, end time.Time) (*models.AcademicTerm, error)
	GetAllTerms() ([]models.AcademicTerm, error)
	GetTerm(id uint) (*models.AcademicTerm, error)
	UpdateTerm(id uint, name string, start, end time.Time, isActive bool) (*models.AcademicTerm, error)
	DeleteTerm(id uint) error

	// Courses
	CreateCourse(code, name string, credits int, desc string, deptID uint, prerequisiteIDs []uint) (*models.Course, error)
	ListCourses(deptID *uint) ([]models.Course, error)
	GetCoursesByDept(deptID uint) ([]models.Course, error)
	GetCourse(id uint) (*models.Course, error)
	UpdateCourse(id uint, code, name string, credits int, desc string, deptID *uint, prerequisiteIDs []uint) (*models.Course, error)
	DeleteCourse(id uint) error

	// Sections
	CreateSection(courseID, termID uint, capacity int, schedule string, professorID *uint) (*models.Section, error)
	ListSections(courseID, termID, profID *uint) ([]models.Section, error)
	GetSectionsByCourse(courseID uint) ([]models.Section, error)
	GetSectionsByTerm(termID uint) ([]models.Section, error)
	GetSectionsByProfessor(profID uint) ([]models.Section, error)
	GetSection(id uint) (*models.Section, error)
	UpdateSection(id uint, courseID, termID *uint, capacity int, schedule string, professorID *uint) (*models.Section, error)
	DeleteSection(id uint) error

	// Teaching assistants
	ListTeachingAssistants(sectionID uint) ([]models.SectionTeachingAssistant, error)
	AssignTeachingAssistant(sectionID, professorID uint, assignedByUserID *uint) (*models.SectionTeachingAssistant, bool, error)
	RemoveTeachingAssistant(sectionID, assignmentID uint) error

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

	ImportGradesFromExcel(file io.Reader, sectionID uint) (int, error)

	// Exam
	CreateExam(sectionID uint, date time.Time, location string) (*models.Exam, error)
	GetExamsBySection(sectionID uint) ([]models.Exam, error)

	// GPA Calculation
	CalculateGPA(studentID uint) (float64, error)
}

type academicService struct {
	repo     repository.AcademicRepository
	userRepo repository.UserRepository
	notifSvc NotificationService
}

func NewAcademicService(repo repository.AcademicRepository, userRepo repository.UserRepository, notifSvc NotificationService) AcademicService {
	return &academicService{repo: repo, userRepo: userRepo, notifSvc: notifSvc}
}

// Term logic
func (s *academicService) CreateTerm(name string, start, end time.Time) (*models.AcademicTerm, error) {
	if end.Before(start) {
		return nil, errors.New("term end date must be after start date")
	}

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
	if end.Before(start) {
		return nil, errors.New("term end date must be after start date")
	}

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

func (s *academicService) ListCourses(deptID *uint) ([]models.Course, error) {
	if deptID != nil {
		return s.repo.GetCoursesByDepartment(*deptID)
	}

	return s.repo.GetAllCourses()
}

// Course logic
func (s *academicService) CreateCourse(code, name string, credits int, desc string, deptID uint, prerequisiteIDs []uint) (*models.Course, error) {
	if credits <= 0 {
		return nil, errors.New("course credits must be greater than zero")
	}
	if s.hasDuplicatePrerequisites(prerequisiteIDs, 0) {
		return nil, errors.New("invalid prerequisite configuration")
	}

	course := &models.Course{Code: code, Name: name, Credits: credits, Description: desc, DepartmentID: deptID}
	if err := s.repo.CreateCourse(course); err != nil {
		return nil, err
	}

	if s.hasDuplicatePrerequisites(prerequisiteIDs, course.ID) {
		return nil, errors.New("invalid prerequisite configuration")
	}
	if err := s.repo.SetCoursePrerequisites(course.ID, prerequisiteIDs); err != nil {
		return nil, err
	}

	return s.repo.GetCourseByID(course.ID)
}

func (s *academicService) hasDuplicatePrerequisites(prerequisiteIDs []uint, courseID uint) bool {
	seen := make(map[uint]struct{}, len(prerequisiteIDs))
	for _, prerequisiteID := range prerequisiteIDs {
		if prerequisiteID == courseID {
			return true
		}
		if _, exists := seen[prerequisiteID]; exists {
			return true
		}
		seen[prerequisiteID] = struct{}{}
	}

	return false
}

func (s *academicService) validateProfessorAssignment(professorID *uint) error {
	if professorID == nil {
		return nil
	}

	professor, err := s.userRepo.FindByID(*professorID)
	if err != nil {
		return errors.New("assigned professor not found")
	}
	if professor.Role != models.RoleProfessor {
		return errors.New("assigned user is not a professor")
	}

	return nil
}

func (s *academicService) validateEnrollmentPrerequisites(studentID uint, section *models.Section) error {
	if section.Course == nil || len(section.Course.Prerequisites) == 0 {
		return nil
	}

	enrollments, err := s.repo.GetEnrollmentsByStudent(studentID)
	if err != nil {
		return err
	}

	completedCourses := make(map[uint]bool)
	for _, enrollment := range enrollments {
		if enrollment.Section == nil || enrollment.Section.Course == nil || enrollment.Grade == nil {
			continue
		}

		if *enrollment.Grade >= 60 {
			completedCourses[enrollment.Section.Course.ID] = true
		}
	}

	for _, prerequisite := range section.Course.Prerequisites {
		if !completedCourses[prerequisite.ID] {
			return fmt.Errorf("missing prerequisite %s", prerequisite.Code)
		}
	}

	return nil
}

func (s *academicService) ensureNoDuplicateEnrollment(studentID, sectionID uint) error {
	_, err := s.repo.GetEnrollment(studentID, sectionID)
	if err == nil {
		return errors.New("student is already enrolled in this section")
	}

	return nil
}

func (s *academicService) UpdateCourse(id uint, code, name string, credits int, desc string, deptID *uint, prerequisiteIDs []uint) (*models.Course, error) {
	if credits <= 0 {
		return nil, errors.New("course credits must be greater than zero")
	}
	if s.hasDuplicatePrerequisites(prerequisiteIDs, id) {
		return nil, errors.New("invalid prerequisite configuration")
	}

	course, err := s.repo.GetCourseByID(id)
	if err != nil {
		return nil, err
	}
	course.Code = code
	course.Name = name
	course.Credits = credits
	course.Description = desc
	if deptID != nil {
		course.DepartmentID = *deptID
	}
	if err := s.repo.UpdateCourse(course); err != nil {
		return nil, err
	}
	if err := s.repo.SetCoursePrerequisites(course.ID, prerequisiteIDs); err != nil {
		return nil, err
	}
	return s.repo.GetCourseByID(course.ID)
}

func (s *academicService) CreateSection(courseID, termID uint, capacity int, schedule string, professorID *uint) (*models.Section, error) {
	if capacity <= 0 {
		return nil, errors.New("section capacity must be greater than zero")
	}
	if err := s.validateProfessorAssignment(professorID); err != nil {
		return nil, err
	}

	section := &models.Section{
		CourseID:       courseID,
		AcademicTermID: termID,
		Capacity:       capacity,
		Schedule:       schedule,
		ProfessorID:    professorID,
	}
	if err := s.repo.CreateSection(section); err != nil {
		return nil, err
	}
	return s.repo.GetSectionByID(section.ID)
}

func (s *academicService) ListSections(courseID, termID, profID *uint) ([]models.Section, error) {
	switch {
	case courseID != nil:
		return s.repo.GetSectionsByCourse(*courseID)
	case termID != nil:
		return s.repo.GetSectionsByTerm(*termID)
	case profID != nil:
		return s.repo.GetSectionsByProfessor(*profID)
	default:
		return s.repo.GetAllSections()
	}
}

func (s *academicService) UpdateSection(id uint, courseID, termID *uint, capacity int, schedule string, professorID *uint) (*models.Section, error) {
	if capacity <= 0 {
		return nil, errors.New("section capacity must be greater than zero")
	}
	if err := s.validateProfessorAssignment(professorID); err != nil {
		return nil, err
	}

	section, err := s.repo.GetSectionByID(id)
	if err != nil {
		return nil, err
	}
	if courseID != nil {
		section.CourseID = *courseID
	}
	if termID != nil {
		section.AcademicTermID = *termID
	}
	section.Capacity = capacity
	section.Schedule = schedule
	section.ProfessorID = professorID
	if err := s.repo.UpdateSection(section); err != nil {
		return nil, err
	}
	return s.repo.GetSectionByID(section.ID)
}

func (s *academicService) EnrollStudent(studentID, sectionID uint) (*models.Enrollment, error) {
	student, err := s.userRepo.FindByID(studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}
	if student.Role != models.RoleStudent {
		return nil, errors.New("user is not a student")
	}

	if err := s.ensureNoDuplicateEnrollment(studentID, sectionID); err != nil {
		return nil, err
	}

	section, err := s.repo.GetSectionByID(sectionID)
	if err != nil {
		return nil, errors.New("section not found")
	}

	if err := s.validateEnrollmentPrerequisites(studentID, section); err != nil {
		return nil, err
	}

	enrollments, _ := s.repo.GetEnrollmentsBySection(sectionID)
	if int64(len(enrollments)) >= int64(section.Capacity) {
		return nil, errors.New("section is at full capacity")
	}

	enrollment := &models.Enrollment{StudentID: studentID, SectionID: sectionID, Status: "enrolled"}
	if err := s.repo.CreateEnrollment(enrollment); err != nil {
		return nil, err
	}
	return s.repo.GetEnrollment(studentID, sectionID)
}

func (s *academicService) UpdateGrade(studentID, sectionID uint, grade float64) (*models.Enrollment, error) {
	if grade < 0 || grade > 100 {
		return nil, errors.New("grade must be between 0 and 100")
	}

	enrollment, err := s.repo.GetEnrollment(studentID, sectionID)
	if err != nil {
		return nil, err
	}
	enrollment.Grade = &grade
	if err := s.repo.UpdateEnrollment(enrollment); err != nil {
		return nil, err
	}
	return s.repo.GetEnrollment(studentID, sectionID)
}

func (s *academicService) DeleteCourse(id uint) error {
	return s.repo.DeleteCourse(id)
}

func (s *academicService) GetCoursesByDept(deptID uint) ([]models.Course, error) {
	return s.repo.GetCoursesByDepartment(deptID)
}

func (s *academicService) GetCourse(id uint) (*models.Course, error) {
	return s.repo.GetCourseByID(id)
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

func (s *academicService) DeleteSection(id uint) error {
	return s.repo.DeleteSection(id)
}

func (s *academicService) ListTeachingAssistants(sectionID uint) ([]models.SectionTeachingAssistant, error) {
	if _, err := s.repo.GetSectionByID(sectionID); err != nil {
		return nil, errors.New("section not found")
	}

	return s.repo.ListTeachingAssistants(sectionID)
}

func (s *academicService) AssignTeachingAssistant(sectionID, professorID uint, assignedByUserID *uint) (*models.SectionTeachingAssistant, bool, error) {
	if _, err := s.repo.GetSectionByID(sectionID); err != nil {
		return nil, false, errors.New("section not found")
	}

	if err := s.validateProfessorAssignment(&professorID); err != nil {
		return nil, false, err
	}

	existing, err := s.repo.GetTeachingAssistantByProfessor(sectionID, professorID)
	if err == nil {
		return existing, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	assignment := &models.SectionTeachingAssistant{
		SectionID:        sectionID,
		ProfessorID:      professorID,
		AssignedByUserID: assignedByUserID,
	}
	if err := s.repo.CreateTeachingAssistant(assignment); err != nil {
		return nil, false, err
	}

	createdAssignment, err := s.repo.GetTeachingAssistant(sectionID, assignment.ID)
	return createdAssignment, true, err
}

func (s *academicService) RemoveTeachingAssistant(sectionID, assignmentID uint) error {
	assignment, err := s.repo.GetTeachingAssistant(sectionID, assignmentID)
	if err != nil {
		return errors.New("teaching assistant assignment not found")
	}

	return s.repo.DeleteTeachingAssistant(assignment.ID)
}

func (s *academicService) GetStudentEnrollments(studentID uint) ([]models.Enrollment, error) {
	return s.repo.GetEnrollmentsByStudent(studentID)
}

func (s *academicService) GetSectionEnrollments(sectionID uint) ([]models.Enrollment, error) {
	return s.repo.GetEnrollmentsBySection(sectionID)
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

func (s *academicService) ImportGradesFromExcel(file io.Reader, sectionID uint) (int, error) {
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
	// Assume first row is header: Student ID, Grade
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}

		studentID, err := strconv.ParseUint(row[0], 10, 32)
		if err != nil {
			continue
		}

		grade, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			continue
		}

		if _, err := s.UpdateGrade(uint(studentID), sectionID, grade); err == nil {
			successCount++
		}
	}

	if successCount > 0 {
		section, _ := s.repo.GetSectionByID(sectionID)
		courseName := ""
		if section != nil && section.Course != nil {
			courseName = section.Course.Name
		}
		s.notifSvc.NotifyStudentsInSections([]uint{sectionID}, "Final Grades Published", fmt.Sprintf("Final grades for %s have been published.", courseName))
	}

	return successCount, nil
}

func (s *academicService) CreateExam(sectionID uint, date time.Time, location string) (*models.Exam, error) {
	exam := &models.Exam{SectionID: sectionID, Date: date, Location: location}
	if err := s.repo.CreateExam(exam); err != nil {
		return nil, err
	}

	// Trigger notification
	section, _ := s.repo.GetSectionByID(sectionID)
	courseName := ""
	if section != nil && section.Course != nil {
		courseName = section.Course.Name
	}
	s.notifSvc.NotifyStudentsInSections([]uint{sectionID}, "Exam Scheduled", fmt.Sprintf("An exam for %s has been scheduled for %s at %s.", courseName, date.Format("Jan 02, 2006 15:04"), location))

	return exam, nil
}

func (s *academicService) GetExamsBySection(sectionID uint) ([]models.Exam, error) {
	return s.repo.GetExamsBySection(sectionID)
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
