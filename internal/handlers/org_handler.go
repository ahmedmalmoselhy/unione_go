package handlers

import (
	"net/http"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type OrgHandler struct {
	orgService services.OrgService
}

func NewOrgHandler(orgService services.OrgService) *OrgHandler {
	return &OrgHandler{orgService: orgService}
}

type CreateOrgInput struct {
	Name string `json:"name" binding:"required"`
}

type AdminCreateFacultyInput struct {
	Name         string `json:"name" binding:"required"`
	UniversityID uint   `json:"university_id" binding:"required"`
}

type AdminUpdateFacultyInput struct {
	Name         string `json:"name" binding:"required"`
	UniversityID *uint  `json:"university_id"`
}

type AdminCreateDepartmentInput struct {
	Name      string `json:"name" binding:"required"`
	FacultyID uint   `json:"faculty_id" binding:"required"`
}

type AdminUpdateDepartmentInput struct {
	Name      string `json:"name" binding:"required"`
	FacultyID *uint  `json:"faculty_id"`
}

func (h *OrgHandler) CreateUniversity(c *gin.Context) {
	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	uni, err := h.orgService.CreateUniversity(input.Name)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "create_university_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, uni)
}

func (h *OrgHandler) GetUniversities(c *gin.Context) {
	unis, err := h.orgService.GetUniversities()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_universities_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, unis)
}

func (h *OrgHandler) UpdateUniversity(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	uni, err := h.orgService.UpdateUniversity(id, input.Name)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "update_university_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, uni)
}

func (h *OrgHandler) DeleteUniversity(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	if err := h.orgService.DeleteUniversity(id); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "delete_university_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *OrgHandler) CreateFaculty(c *gin.Context) {
	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	universityID, err := apiutil.ParseUintParam(c, "university_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_university_id", err.Error())
		return
	}

	faculty, err := h.orgService.CreateFaculty(input.Name, universityID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "create_faculty_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, faculty)
}

func (h *OrgHandler) GetFaculties(c *gin.Context) {
	universityID, err := apiutil.ParseUintParam(c, "university_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_university_id", err.Error())
		return
	}

	faculties, err := h.orgService.GetFaculties(universityID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_faculties_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, faculties)
}

func (h *OrgHandler) UpdateFaculty(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	faculty, err := h.orgService.UpdateFaculty(id, input.Name, nil)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "update_faculty_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, faculty)
}

func (h *OrgHandler) DeleteFaculty(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	if err := h.orgService.DeleteFaculty(id); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "delete_faculty_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *OrgHandler) CreateDepartment(c *gin.Context) {
	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	facultyID, err := apiutil.ParseUintParam(c, "faculty_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_faculty_id", err.Error())
		return
	}

	department, err := h.orgService.CreateDepartment(input.Name, facultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "create_department_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, department)
}

func (h *OrgHandler) GetDepartments(c *gin.Context) {
	facultyID, err := apiutil.ParseUintParam(c, "faculty_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_faculty_id", err.Error())
		return
	}

	departments, err := h.orgService.GetDepartments(facultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_departments_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, departments)
}

func (h *OrgHandler) UpdateDepartment(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	department, err := h.orgService.UpdateDepartment(id, input.Name, nil)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "update_department_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, department)
}

func (h *OrgHandler) DeleteDepartment(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	if err := h.orgService.DeleteDepartment(id); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "delete_department_failed", err.Error())
		return
	}

	apiutil.NoContent(c)
}

func (h *OrgHandler) ListFaculties(c *gin.Context) {
	universityID, err := apiutil.ParseOptionalUintQuery(c, "university_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_university_id", err.Error())
		return
	}

	faculties, err := h.orgService.ListFaculties(universityID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_faculties_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, faculties)
}

func (h *OrgHandler) AdminCreateFaculty(c *gin.Context) {
	var input AdminCreateFacultyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	faculty, err := h.orgService.CreateFaculty(input.Name, input.UniversityID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "create_faculty_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, faculty)
}

func (h *OrgHandler) AdminUpdateFaculty(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	var input AdminUpdateFacultyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	faculty, err := h.orgService.UpdateFaculty(id, input.Name, input.UniversityID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "update_faculty_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, faculty)
}

func (h *OrgHandler) ListDepartments(c *gin.Context) {
	facultyID, err := apiutil.ParseOptionalUintQuery(c, "faculty_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_faculty_id", err.Error())
		return
	}

	departments, err := h.orgService.ListDepartments(facultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "list_departments_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, departments)
}

func (h *OrgHandler) AdminCreateDepartment(c *gin.Context) {
	var input AdminCreateDepartmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	department, err := h.orgService.CreateDepartment(input.Name, input.FacultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "create_department_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, department)
}

func (h *OrgHandler) AdminUpdateDepartment(c *gin.Context) {
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	var input AdminUpdateDepartmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	department, err := h.orgService.UpdateDepartment(id, input.Name, input.FacultyID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "update_department_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, department)
}
