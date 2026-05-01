package handlers

import (
	"net/http"
	"strconv"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentService services.StudentService
}

func NewStudentHandler(studentService services.StudentService) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}

type CreateStudentInput struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	FacultyID    uint   `json:"faculty_id" binding:"required"`
	DepartmentID *uint  `json:"department_id"`
}

type UpdateStudentInput struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	FacultyID    uint   `json:"faculty_id" binding:"required"`
	DepartmentID *uint  `json:"department_id"`
}

type TransferStudentInput struct {
	FacultyID    uint  `json:"faculty_id" binding:"required"`
	DepartmentID *uint `json:"department_id"`
}

// CreateStudent godoc
// @Summary Create a new student
// @Description Creates a new student user and profile.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body CreateStudentInput true "Student details"
// @Success 201 {object} models.User
// @Failure 400 {object} apiutil.errorResponse
// @Router /admin/students [post]
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var input CreateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	student, err := h.studentService.CreateStudent(
		input.Email,
		input.Password,
		input.FirstName,
		input.LastName,
		input.FacultyID,
		input.DepartmentID,
	)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "creation_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, student)
}

// ListStudents godoc
// @Summary List all students
// @Description Returns a list of all students, optionally filtered by faculty or department.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param faculty_id query int false "Filter by faculty ID"
// @Param department_id query int false "Filter by department ID"
// @Success 200 {array} models.User
// @Router /admin/students [get]
func (h *StudentHandler) ListStudents(c *gin.Context) {
	var facultyID, departmentID *uint

	if fIDStr := c.Query("faculty_id"); fIDStr != "" {
		if id, err := strconv.ParseUint(fIDStr, 10, 32); err == nil {
			val := uint(id)
			facultyID = &val
		}
	}

	if dIDStr := c.Query("department_id"); dIDStr != "" {
		if id, err := strconv.ParseUint(dIDStr, 10, 32); err == nil {
			val := uint(id)
			departmentID = &val
		}
	}

	students, err := h.studentService.ListStudents(facultyID, departmentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, students)
}

// GetStudent godoc
// @Summary Get student by ID
// @Description Returns a specific student by their ID.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Success 200 {object} models.User
// @Failure 404 {object} apiutil.errorResponse
// @Router /admin/students/{id} [get]
func (h *StudentHandler) GetStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	student, err := h.studentService.GetStudent(uint(id))
	if err != nil {
		apiutil.Error(c, http.StatusNotFound, "not_found", "Student not found")
		return
	}

	apiutil.Success(c, http.StatusOK, student)
}

// UpdateStudent godoc
// @Summary Update student details
// @Description Updates the specified student's information.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Param input body UpdateStudentInput true "Updated details"
// @Success 200 {object} models.User
// @Failure 400 {object} apiutil.errorResponse
// @Router /admin/students/{id} [put]
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var input UpdateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	student, err := h.studentService.UpdateStudent(
		uint(id),
		input.FirstName,
		input.LastName,
		input.FacultyID,
		input.DepartmentID,
	)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "update_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, student)
}

// DeleteStudent godoc
// @Summary Delete a student
// @Description Soft-deletes a student by their ID.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Success 204 "No Content"
// @Router /admin/students/{id} [delete]
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.studentService.DeleteStudent(uint(id)); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "delete_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

// TransferStudent godoc
// @Summary Transfer a student to another faculty/department
// @Description Record a student transfer and update their current organization link.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Param input body TransferStudentInput true "Transfer details"
// @Success 200 {object} models.User
// @Router /admin/students/{id}/transfer [post]
func (h *StudentHandler) TransferStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var input TransferStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	student, err := h.studentService.TransferStudent(uint(id), input.FacultyID, input.DepartmentID)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "transfer_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, student)
}

// GetTransferHistory godoc
// @Summary Get student transfer history
// @Description Returns the history of organization transfers for a student.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Success 200 {array} models.StudentDepartmentHistory
// @Router /admin/students/{id}/transfers [get]
func (h *StudentHandler) GetTransferHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	history, err := h.studentService.GetTransferHistory(uint(id))
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "history_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, history)
}

// ExportStudents godoc
// @Summary Export students to Excel
// @Description Returns an Excel file containing the list of students, optionally filtered.
// @Tags admin
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Param faculty_id query int false "Filter by faculty ID"
// @Param department_id query int false "Filter by department ID"
// @Success 200 {string} string "Excel binary data"
// @Router /admin/students/export [get]
func (h *StudentHandler) ExportStudents(c *gin.Context) {
	var facultyID, departmentID *uint

	if fIDStr := c.Query("faculty_id"); fIDStr != "" {
		if id, err := strconv.ParseUint(fIDStr, 10, 32); err == nil {
			val := uint(id)
			facultyID = &val
		}
	}

	if dIDStr := c.Query("department_id"); dIDStr != "" {
		if id, err := strconv.ParseUint(dIDStr, 10, 32); err == nil {
			val := uint(id)
			departmentID = &val
		}
	}

	data, err := h.studentService.ExportStudents(facultyID, departmentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "export_failed", err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename=students.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
