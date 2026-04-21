package handlers

import (
	"net/http"
	"strconv"

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

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	facultyIDStr := c.Param("faculty_id")
	facultyID, err := strconv.ParseUint(facultyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	var input CreateEmployeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.employeeService.CreateEmployee(input.Email, input.Password, input.FirstName, input.LastName, uint(facultyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	facultyIDStr := c.Param("faculty_id")
	facultyID, err := strconv.ParseUint(facultyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	employees, err := h.employeeService.GetEmployeesByFaculty(uint(facultyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var input UpdateEmployeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.employeeService.UpdateEmployee(uint(id), input.FirstName, input.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) ImportStudents(c *gin.Context) {
	facultyIDStr := c.Param("faculty_id")
	facultyID, err := strconv.ParseUint(facultyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer openedFile.Close()

	count, err := h.employeeService.ImportStudentsFromExcel(openedFile, uint(facultyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"imported_count": count})
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := h.employeeService.DeleteEmployee(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
