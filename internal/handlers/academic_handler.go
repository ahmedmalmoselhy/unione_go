package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type AcademicHandler struct {
	academicService services.AcademicService
}

func NewAcademicHandler(academicService services.AcademicService) *AcademicHandler {
	return &AcademicHandler{academicService: academicService}
}

// Terms
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

func (h *AcademicHandler) CreateTerm(c *gin.Context) {
	var input CreateTermInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	term, err := h.academicService.CreateTerm(input.Name, input.StartDate, input.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, term)
}

func (h *AcademicHandler) GetTerms(c *gin.Context) {
	terms, err := h.academicService.GetAllTerms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, terms)
}

// Courses
type CreateCourseInput struct {
	Code         string `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Credits      int    `json:"credits" binding:"required"`
	Description  string `json:"description"`
	DepartmentID uint   `json:"department_id" binding:"required"`
}

func (h *AcademicHandler) CreateCourse(c *gin.Context) {
	var input CreateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	course, err := h.academicService.CreateCourse(input.Code, input.Name, input.Credits, input.Description, input.DepartmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)
}

func (h *AcademicHandler) GetCourses(c *gin.Context) {
	deptIDStr := c.Query("department_id")
	if deptIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "department_id is required"})
		return
	}
	deptID, _ := strconv.ParseUint(deptIDStr, 10, 32)
	courses, err := h.academicService.GetCoursesByDept(uint(deptID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

// Sections
type CreateSectionInput struct {
	CourseID       uint   `json:"course_id" binding:"required"`
	AcademicTermID uint   `json:"academic_term_id" binding:"required"`
	Capacity       int    `json:"capacity" binding:"required"`
	Schedule       string `json:"schedule"`
}

func (h *AcademicHandler) CreateSection(c *gin.Context) {
	var input CreateSectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	section, err := h.academicService.CreateSection(input.CourseID, input.AcademicTermID, input.Capacity, input.Schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, section)
}

func (h *AcademicHandler) GetSections(c *gin.Context) {
	courseIDStr := c.Query("course_id")
	termIDStr := c.Query("academic_term_id")
	profIDStr := c.Query("professor_id")
	
	if courseIDStr != "" {
		courseID, _ := strconv.ParseUint(courseIDStr, 10, 32)
		sections, err := h.academicService.GetSectionsByCourse(uint(courseID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sections)
		return
	}
	
	if termIDStr != "" {
		termID, _ := strconv.ParseUint(termIDStr, 10, 32)
		sections, err := h.academicService.GetSectionsByTerm(uint(termID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sections)
		return
	}

	if profIDStr != "" {
		profID, _ := strconv.ParseUint(profIDStr, 10, 32)
		sections, err := h.academicService.GetSectionsByProfessor(uint(profID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sections)
		return
	}
	
	c.JSON(http.StatusBadRequest, gin.H{"error": "course_id, academic_term_id or professor_id is required"})
}

// Enrollments
func (h *AcademicHandler) Enroll(c *gin.Context) {
	sectionIDStr := c.Param("section_id")
	sectionID, _ := strconv.ParseUint(sectionIDStr, 10, 32)
	
	var input struct {
		StudentID uint `json:"student_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enrollment, err := h.academicService.EnrollStudent(input.StudentID, uint(sectionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, enrollment)
}

func (h *AcademicHandler) GetStudentEnrollments(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	studentID, _ := strconv.ParseUint(studentIDStr, 10, 32)
	
	enrollments, err := h.academicService.GetStudentEnrollments(uint(studentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enrollments)
}

type UpdateGradeInput struct {
	Grade float64 `json:"grade" binding:"required"`
}

func (h *AcademicHandler) UpdateGrade(c *gin.Context) {
	sectionIDStr := c.Param("section_id")
	studentIDStr := c.Param("student_id")
	sectionID, _ := strconv.ParseUint(sectionIDStr, 10, 32)
	studentID, _ := strconv.ParseUint(studentIDStr, 10, 32)

	var input UpdateGradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enrollment, err := h.academicService.UpdateGrade(uint(studentID), uint(sectionID), input.Grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enrollment)
}

// Attendance
type RecordAttendanceInput struct {
	Date            time.Time    `json:"date" binding:"required"`
	StudentPresence map[uint]bool `json:"student_presence" binding:"required"`
}

func (h *AcademicHandler) RecordAttendance(c *gin.Context) {
	sectionIDStr := c.Param("section_id")
	sectionID, _ := strconv.ParseUint(sectionIDStr, 10, 32)

	var input RecordAttendanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.academicService.RecordAttendance(uint(sectionID), input.Date, input.StudentPresence); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance recorded successfully"})
}

func (h *AcademicHandler) GetSectionAttendance(c *gin.Context) {
	sectionIDStr := c.Param("section_id")
	dateStr := c.Query("date")
	sectionID, _ := strconv.ParseUint(sectionIDStr, 10, 32)
	
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use RFC3339"})
		return
	}

	attendance, err := h.academicService.GetSectionAttendance(uint(sectionID), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attendance)
}

// GPA
func (h *AcademicHandler) GetGPA(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	studentID, _ := strconv.ParseUint(studentIDStr, 10, 32)

	gpa, err := h.academicService.CalculateGPA(uint(studentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"student_id": studentID, "gpa": gpa})
}
