package handlers

import (
	"net/http"
	"strconv"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	employeeService services.EmployeeService
}

func NewEmployeeHandler(employeeService services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
}

type CreateEmployeeInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type UpdateEmployeeInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type AdminCreateEmployeeInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	FacultyID uint   `json:"faculty_id" binding:"required"`
}

type AdminUpdateEmployeeInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	FacultyID uint   `json:"faculty_id" binding:"required"`
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	facultyID, err := apiutil.ParseUintParam(c, "faculty_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_faculty_id", err.Error())
		return
	}

	var input CreateEmployeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	employee, err := h.employeeService.CreateEmployee(input.Email, input.Password, input.FirstName, input.LastName, facultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "create_employee_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, employee)
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	facultyID, err := apiutil.ParseUintParam(c, "faculty_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_faculty_id", err.Error())
		return
	}

	employees, err := h.employeeService.GetEmployeesByFaculty(facultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_employees_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, employees)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_employee_id", err.Error())
		return
	}

	var input UpdateEmployeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	employee, err := h.employeeService.UpdateEmployee(id, input.FirstName, input.LastName)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "update_employee_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, employee)
}

func (h *EmployeeHandler) ImportStudents(c *gin.Context) {
	facultyID, err := apiutil.ParseUintParam(c, "faculty_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_faculty_id", err.Error())
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

	count, err := h.employeeService.ImportStudentsFromExcel(openedFile, facultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "student_import_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{"imported_count": count})
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_employee_id", err.Error())
		return
	}

	if err := h.employeeService.DeleteEmployee(id); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "delete_employee_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	var facultyID *uint
	if c.Query("faculty_id") != "" {
		value, err := apiutil.ParseRequiredUintQuery(c, "faculty_id")
		if err != nil {
			apiutil.Error(c, http.StatusBadRequest, "invalid_faculty_id", err.Error())
			return
		}
		facultyID = &value
	}

	employees, err := h.employeeService.ListEmployees(facultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_employees_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, employees)
}

func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_employee_id", err.Error())
		return
	}

	employee, err := h.employeeService.GetEmployee(id)
	if err != nil {
		apiutil.Error(c, http.StatusNotFound, "employee_not_found", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, employee)
}

func (h *EmployeeHandler) AdminCreateEmployee(c *gin.Context) {
	var input AdminCreateEmployeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	employee, err := h.employeeService.CreateEmployee(input.Email, input.Password, input.FirstName, input.LastName, input.FacultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "create_employee_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, employee)
}

func (h *EmployeeHandler) AdminUpdateEmployee(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_employee_id", err.Error())
		return
	}

	var input AdminUpdateEmployeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	employee, err := h.employeeService.UpdateEmployeeDetails(id, input.FirstName, input.LastName, input.FacultyID)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "update_employee_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, employee)
}

func (h *EmployeeHandler) ExportEmployees(c *gin.Context) {
	var facultyID *uint
	if fID := c.Query("faculty_id"); fID != "" {
		id, err := strconv.ParseUint(fID, 10, 32)
		if err == nil {
			uID := uint(id)
			facultyID = &uID
		}
	}

	data, err := h.employeeService.ExportEmployees(facultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "export_failed", err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=employees.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
