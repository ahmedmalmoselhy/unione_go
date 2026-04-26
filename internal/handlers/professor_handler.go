package handlers

import (
	"net/http"
	"strconv"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type ProfessorHandler struct {
	professorService services.ProfessorService
}

func NewProfessorHandler(professorService services.ProfessorService) *ProfessorHandler {
	return &ProfessorHandler{professorService: professorService}
}

type CreateProfessorInput struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	FacultyID    uint   `json:"faculty_id" binding:"required"`
	DepartmentID *uint  `json:"department_id"`
}

type UpdateProfessorInput struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	FacultyID    uint   `json:"faculty_id" binding:"required"`
	DepartmentID *uint  `json:"department_id"`
}

func (h *ProfessorHandler) CreateProfessor(c *gin.Context) {
	var input CreateProfessorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	professor, err := h.professorService.CreateProfessor(
		input.Email,
		input.Password,
		input.FirstName,
		input.LastName,
		input.FacultyID,
		input.DepartmentID,
	)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "create_professor_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, professor)
}

func (h *ProfessorHandler) ListProfessors(c *gin.Context) {
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

	professors, err := h.professorService.ListProfessors(facultyID, departmentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_professors_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, professors)
}

func (h *ProfessorHandler) GetProfessor(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_professor_id", err.Error())
		return
	}

	professor, err := h.professorService.GetProfessor(id)
	if err != nil {
		apiutil.Error(c, http.StatusNotFound, "professor_not_found", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, professor)
}

func (h *ProfessorHandler) UpdateProfessor(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_professor_id", err.Error())
		return
	}

	var input UpdateProfessorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	professor, err := h.professorService.UpdateProfessor(
		id,
		input.FirstName,
		input.LastName,
		input.FacultyID,
		input.DepartmentID,
	)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "update_professor_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, professor)
}

func (h *ProfessorHandler) DeleteProfessor(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_professor_id", err.Error())
		return
	}

	if err := h.professorService.DeleteProfessor(id); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "delete_professor_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *ProfessorHandler) ExportProfessors(c *gin.Context) {
	var facultyID *uint
	if fID := c.Query("faculty_id"); fID != "" {
		id, err := strconv.ParseUint(fID, 10, 32)
		if err == nil {
			uID := uint(id)
			facultyID = &uID
		}
	}

	var deptID *uint
	if dID := c.Query("department_id"); dID != "" {
		id, err := strconv.ParseUint(dID, 10, 32)
		if err == nil {
			uID := uint(id)
			deptID = &uID
		}
	}

	data, err := h.professorService.ExportProfessors(facultyID, deptID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "export_failed", err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=professors.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
