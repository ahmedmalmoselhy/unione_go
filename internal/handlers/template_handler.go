package handlers

import (
	"net/http"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	impExpSvc services.ImportExportService
}

func NewTemplateHandler(impExpSvc services.ImportExportService) *TemplateHandler {
	return &TemplateHandler{impExpSvc: impExpSvc}
}

type StudentImportTemplate struct {
	NationalID   string `excel:"National ID"`
	Email        string `excel:"Email"`
	FirstName    string `excel:"First Name"`
	LastName     string `excel:"Last Name"`
	Phone        string `excel:"Phone"`
	Gender       string `excel:"Gender"`
	DateOfBirth  string `excel:"Date of Birth (YYYY-MM-DD)"`
	DepartmentID string `excel:"Department ID"`
}

func (h *TemplateHandler) DownloadStudentsImportTemplate(c *gin.Context) {
	data := []StudentImportTemplate{
		{
			NationalID:   "123456789",
			Email:        "student@example.edu",
			FirstName:    "Jane",
			LastName:     "Doe",
			Phone:        "0123456789",
			Gender:       "female",
			DateOfBirth:  "2000-01-01",
			DepartmentID: "1",
		},
	}

	excelData, err := h.impExpSvc.ExportToExcel(data, "Students Template")
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "template_generation_failed", err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=students_import_template.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}

type GradeImportTemplate struct {
	StudentID string `excel:"Student ID or Email"`
	Grade     string `excel:"Grade"`
}

func (h *TemplateHandler) DownloadGradesImportTemplate(c *gin.Context) {
	data := []GradeImportTemplate{
		{
			StudentID: "1001",
			Grade:     "92.5",
		},
	}

	excelData, err := h.impExpSvc.ExportToExcel(data, "Grades Template")
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "template_generation_failed", err.Error())
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=grades_import_template.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}
