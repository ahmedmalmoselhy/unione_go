package handlers

import (
	"net/http"

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
		apiutil.Error(c, http.StatusBadRequest, "create_student_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, student)
}

func (h *StudentHandler) ListStudents(c *gin.Context) {
	var facultyID *uint
	if c.Query("faculty_id") != "" {
		value, err := apiutil.ParseRequiredUintQuery(c, "faculty_id")
		if err != nil {
			apiutil.Error(c, http.StatusBadRequest, "invalid_faculty_id", err.Error())
			return
		}
		facultyID = &value
	}

	var departmentID *uint
	if c.Query("department_id") != "" {
		value, err := apiutil.ParseRequiredUintQuery(c, "department_id")
		if err != nil {
			apiutil.Error(c, http.StatusBadRequest, "invalid_department_id", err.Error())
			return
		}
		departmentID = &value
	}

	students, err := h.studentService.ListStudents(facultyID, departmentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_students_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, students)
}

func (h *StudentHandler) GetStudent(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_student_id", err.Error())
		return
	}

	student, err := h.studentService.GetStudent(id)
	if err != nil {
		apiutil.Error(c, http.StatusNotFound, "student_not_found", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_student_id", err.Error())
		return
	}

	var input UpdateStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	student, err := h.studentService.UpdateStudent(
		id,
		input.FirstName,
		input.LastName,
		input.FacultyID,
		input.DepartmentID,
	)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "update_student_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, student)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_student_id", err.Error())
		return
	}

	if err := h.studentService.DeleteStudent(id); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "delete_student_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *StudentHandler) TransferStudent(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_student_id", err.Error())
		return
	}

	var input TransferStudentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	student, err := h.studentService.TransferStudent(id, input.FacultyID, input.DepartmentID)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "transfer_student_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, student)
}

func (h *StudentHandler) GetTransferHistory(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_student_id", err.Error())
		return
	}

	history, err := h.studentService.GetTransferHistory(id)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "transfer_history_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, history)
}
