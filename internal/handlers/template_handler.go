package handlers

import "github.com/gin-gonic/gin"

type TemplateHandler struct{}

func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{}
}

func (h *TemplateHandler) DownloadStudentsImportTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="students_import_template.csv"`)
	c.String(200, "email,password,first_name,last_name,department_id\nstudent@example.edu,ChangeMe123,Jane,Doe,12\n")
}

func (h *TemplateHandler) DownloadGradesImportTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="grades_import_template.csv"`)
	c.String(200, "student_id,grade\n1001,92.5\n")
}
