package handlers

import (
	"net/http"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/middlewares"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type AcademicHandler struct {
	academicService services.AcademicService
}

func NewAcademicHandler(academicService services.AcademicService) *AcademicHandler {
	return &AcademicHandler{academicService: academicService}
}

type CreateTermInput struct {
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type UpdateTermInput struct {
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	IsActive  bool      `json:"is_active"`
}

type CreateCourseInput struct {
	Code            string `json:"code" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Credits         int    `json:"credits" binding:"required"`
	Description     string `json:"description"`
	DepartmentID    uint   `json:"department_id" binding:"required"`
	PrerequisiteIDs []uint `json:"prerequisite_ids"`
}

type UpdateCourseInput struct {
	Code            string `json:"code" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Credits         int    `json:"credits" binding:"required"`
	Description     string `json:"description"`
	DepartmentID    *uint  `json:"department_id"`
	PrerequisiteIDs []uint `json:"prerequisite_ids"`
}

type CreateSectionInput struct {
	CourseID       uint   `json:"course_id" binding:"required"`
	AcademicTermID uint   `json:"academic_term_id" binding:"required"`
	Capacity       int    `json:"capacity" binding:"required"`
	Schedule       string `json:"schedule"`
	ProfessorID    *uint  `json:"professor_id"`
}

type UpdateSectionInput struct {
	CourseID       *uint  `json:"course_id"`
	AcademicTermID *uint  `json:"academic_term_id"`
	Capacity       int    `json:"capacity" binding:"required"`
	Schedule       string `json:"schedule"`
	ProfessorID    *uint  `json:"professor_id"`
}

type EnrollmentInput struct {
	StudentID uint `json:"student_id" binding:"required"`
}

type UpdateGradeInput struct {
	Grade float64 `json:"grade" binding:"required"`
}

type RecordAttendanceInput struct {
	Date            time.Time     `json:"date" binding:"required"`
	StudentPresence map[uint]bool `json:"student_presence" binding:"required"`
}

type CreateExamInput struct {
	Date     time.Time `json:"date" binding:"required"`
	Location string    `json:"location" binding:"required"`
}

type AssignTeachingAssistantInput struct {
	ProfessorID uint `json:"professor_id" binding:"required"`
}

func (h *AcademicHandler) CreateTerm(c *gin.Context) {
	var input CreateTermInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	term, err := h.academicService.CreateTerm(input.Name, input.StartDate, input.EndDate)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "create_term_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, term)
}

func (h *AcademicHandler) GetTerms(c *gin.Context) {
	terms, err := h.academicService.GetAllTerms()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_terms_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, terms)
}

func (h *AcademicHandler) GetTerm(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_term_id", err.Error())
		return
	}

	term, err := h.academicService.GetTerm(id)
	if err != nil {
		apiutil.Error(c, http.StatusNotFound, "term_not_found", "Term not found")
		return
	}

	apiutil.Success(c, http.StatusOK, term)
}

func (h *AcademicHandler) UpdateTerm(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_term_id", err.Error())
		return
	}

	var input UpdateTermInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	term, err := h.academicService.UpdateTerm(id, input.Name, input.StartDate, input.EndDate, input.IsActive)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "update_term_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, term)
}

func (h *AcademicHandler) DeleteTerm(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_term_id", err.Error())
		return
	}

	if err := h.academicService.DeleteTerm(id); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "delete_term_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *AcademicHandler) CreateCourse(c *gin.Context) {
	var input CreateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	course, err := h.academicService.CreateCourse(
		input.Code,
		input.Name,
		input.Credits,
		input.Description,
		input.DepartmentID,
		input.PrerequisiteIDs,
	)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "create_course_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, course)
}

func (h *AcademicHandler) GetCourses(c *gin.Context) {
	departmentID, err := apiutil.ParseOptionalUintQuery(c, "department_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_department_id", err.Error())
		return
	}

	courses, err := h.academicService.ListCourses(departmentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_courses_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, courses)
}

func (h *AcademicHandler) GetCourse(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_course_id", err.Error())
		return
	}

	course, err := h.academicService.GetCourse(id)
	if err != nil {
		apiutil.Error(c, http.StatusNotFound, "course_not_found", "Course not found")
		return
	}

	apiutil.Success(c, http.StatusOK, course)
}

func (h *AcademicHandler) UpdateCourse(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_course_id", err.Error())
		return
	}

	var input UpdateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	course, err := h.academicService.UpdateCourse(id, input.Code, input.Name, input.Credits, input.Description, input.DepartmentID, input.PrerequisiteIDs)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "update_course_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, course)
}

func (h *AcademicHandler) DeleteCourse(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_course_id", err.Error())
		return
	}

	if err := h.academicService.DeleteCourse(id); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "delete_course_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *AcademicHandler) CreateSection(c *gin.Context) {
	var input CreateSectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	section, err := h.academicService.CreateSection(input.CourseID, input.AcademicTermID, input.Capacity, input.Schedule, input.ProfessorID)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "create_section_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, section)
}

func (h *AcademicHandler) GetSections(c *gin.Context) {
	courseID, err := apiutil.ParseOptionalUintQuery(c, "course_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_course_id", err.Error())
		return
	}

	termID, err := apiutil.ParseOptionalUintQuery(c, "academic_term_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_term_id", err.Error())
		return
	}

	professorID, err := apiutil.ParseOptionalUintQuery(c, "professor_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_professor_id", err.Error())
		return
	}

	sections, err := h.academicService.ListSections(courseID, termID, professorID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_sections_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, sections)
}

func (h *AcademicHandler) GetSection(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	section, err := h.academicService.GetSection(id)
	if err != nil {
		apiutil.Error(c, http.StatusNotFound, "section_not_found", "Section not found")
		return
	}

	apiutil.Success(c, http.StatusOK, section)
}

func (h *AcademicHandler) UpdateSection(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	var input UpdateSectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	section, err := h.academicService.UpdateSection(id, input.CourseID, input.AcademicTermID, input.Capacity, input.Schedule, input.ProfessorID)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "update_section_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, section)
}

func (h *AcademicHandler) DeleteSection(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	if err := h.academicService.DeleteSection(id); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "delete_section_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *AcademicHandler) Enroll(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	var input EnrollmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	enrollment, err := h.academicService.EnrollStudent(input.StudentID, sectionID)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "enrollment_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, enrollment)
}

func (h *AcademicHandler) GetStudentEnrollments(c *gin.Context) {
	studentID, err := apiutil.ParseUintParam(c, "student_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_student_id", err.Error())
		return
	}

	enrollments, err := h.academicService.GetStudentEnrollments(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_enrollments_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, enrollments)
}

func (h *AcademicHandler) GetSectionEnrollments(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	enrollments, err := h.academicService.GetSectionEnrollments(sectionID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_enrollments_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, enrollments)
}

func (h *AcademicHandler) UpdateGrade(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	studentID, err := apiutil.ParseUintParam(c, "student_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_student_id", err.Error())
		return
	}

	var input UpdateGradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	enrollment, err := h.academicService.UpdateGrade(studentID, sectionID, input.Grade)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "update_grade_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, enrollment)
}

func (h *AcademicHandler) DropSection(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	studentID, err := apiutil.ParseUintParam(c, "student_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_student_id", err.Error())
		return
	}

	if err := h.academicService.DropSection(studentID, sectionID); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "drop_enrollment_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *AcademicHandler) RecordAttendance(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	var input RecordAttendanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.academicService.RecordAttendance(sectionID, input.Date, input.StudentPresence); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "record_attendance_failed", err.Error())
		return
	}

	apiutil.Message(c, http.StatusOK, "Attendance recorded successfully")
}

func (h *AcademicHandler) GetSectionAttendance(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	dateStr := c.Query("date")
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_date", "Invalid date format. Use RFC3339")
		return
	}

	attendance, err := h.academicService.GetSectionAttendance(sectionID, date)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_attendance_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, attendance)
}

func (h *AcademicHandler) ImportGrades(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", "No file uploaded")
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "file_open_failed", "Failed to open file")
		return
	}
	defer openedFile.Close()

	count, err := h.academicService.ImportGradesFromExcel(openedFile, sectionID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "grade_import_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{"imported_count": count})
}

func (h *AcademicHandler) CreateExam(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	var input CreateExamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	exam, err := h.academicService.CreateExam(sectionID, input.Date, input.Location)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "create_exam_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, exam)
}

func (h *AcademicHandler) GetExams(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	exams, err := h.academicService.GetExamsBySection(sectionID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_exams_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, exams)
}

func (h *AcademicHandler) GetGPA(c *gin.Context) {
	studentID, err := apiutil.ParseUintParam(c, "student_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_student_id", err.Error())
		return
	}

	gpa, err := h.academicService.CalculateGPA(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "gpa_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{"student_id": studentID, "gpa": gpa})
}

func (h *AcademicHandler) ListTeachingAssistants(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	assignments, err := h.academicService.ListTeachingAssistants(sectionID)
	if err != nil {
		apiutil.Error(c, http.StatusNotFound, "list_teaching_assistants_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, assignments)
}

func (h *AcademicHandler) AssignTeachingAssistant(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	var input AssignTeachingAssistantInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	assignedByUserID, ok := middlewares.GetUserID(c)
	var assignedBy *uint
	if ok {
		assignedBy = &assignedByUserID
	}

	assignment, created, err := h.academicService.AssignTeachingAssistant(sectionID, input.ProfessorID, assignedBy)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "assign_teaching_assistant_failed", err.Error())
		return
	}

	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}

	apiutil.Success(c, status, assignment)
}

func (h *AcademicHandler) RemoveTeachingAssistant(c *gin.Context) {
	sectionID, err := apiutil.ParseUintParam(c, "section_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_section_id", err.Error())
		return
	}

	assignmentID, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_assignment_id", err.Error())
		return
	}

	if err := h.academicService.RemoveTeachingAssistant(sectionID, assignmentID); err != nil {
		apiutil.Error(c, http.StatusNotFound, "remove_teaching_assistant_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}
