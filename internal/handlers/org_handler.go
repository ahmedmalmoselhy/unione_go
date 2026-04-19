package handlers

import (
	"net/http"
	"strconv"

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

func (h *OrgHandler) CreateUniversity(c *gin.Context) {
	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uni, err := h.orgService.CreateUniversity(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, uni)
}

func (h *OrgHandler) GetUniversities(c *gin.Context) {
	unis, err := h.orgService.GetUniversities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, unis)
}

func (h *OrgHandler) CreateFaculty(c *gin.Context) {
	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uniIDStr := c.Param("university_id")
	uniID, err := strconv.ParseUint(uniIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	faculty, err := h.orgService.CreateFaculty(input.Name, uint(uniID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, faculty)
}

func (h *OrgHandler) GetFaculties(c *gin.Context) {
	uniIDStr := c.Param("university_id")
	uniID, err := strconv.ParseUint(uniIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid university ID"})
		return
	}

	faculties, err := h.orgService.GetFaculties(uint(uniID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, faculties)
}

func (h *OrgHandler) CreateDepartment(c *gin.Context) {
	var input CreateOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	facultyIDStr := c.Param("faculty_id")
	facultyID, err := strconv.ParseUint(facultyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	dept, err := h.orgService.CreateDepartment(input.Name, uint(facultyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dept)
}

func (h *OrgHandler) GetDepartments(c *gin.Context) {
	facultyIDStr := c.Param("faculty_id")
	facultyID, err := strconv.ParseUint(facultyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid faculty ID"})
		return
	}

	depts, err := h.orgService.GetDepartments(uint(facultyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, depts)
}
