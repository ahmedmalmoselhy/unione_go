package services

import (
	"fmt"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
)

type TranscriptService interface {
	GenerateTranscriptPDF(student *models.User, enrollments []models.Enrollment, gpa float64) ([]byte, error)
}

type transcriptService struct{}

func NewTranscriptService() TranscriptService {
	return &transcriptService{}
}

func (s *transcriptService) GenerateTranscriptPDF(student *models.User, enrollments []models.Enrollment, gpa float64) ([]byte, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	// Header
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("OFFICIAL ACADEMIC TRANSCRIPT", props.Text{
				Size:  16,
				Align: consts.Center,
				Style: consts.Bold,
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("UniOne University Management System", props.Text{
				Size:  10,
				Align: consts.Center,
			})
		})
	})

	m.Line(5)

	// Student Info
	m.Row(30, func() {
		m.Col(6, func() {
			m.Text(fmt.Sprintf("Student Name: %s %s", student.FirstName, student.LastName), props.Text{Size: 10, Style: consts.Bold})
			m.Text(fmt.Sprintf("Student ID: %d", student.ID), props.Text{Size: 10, Top: 5})
			m.Text(fmt.Sprintf("Email: %s", student.Email), props.Text{Size: 10, Top: 10})
		})
		m.Col(6, func() {
			facultyName := "N/A"
			if student.Faculty != nil {
				facultyName = student.Faculty.Name
			}
			m.Text(fmt.Sprintf("Faculty: %s", facultyName), props.Text{Size: 10, Style: consts.Bold})
			m.Text(fmt.Sprintf("Cumulative GPA: %.2f", gpa), props.Text{Size: 10, Top: 5, Style: consts.Bold})
			m.Text(fmt.Sprintf("Date Issued: %s", time.Now().Format("2006-01-02")), props.Text{Size: 10, Top: 10})
		})
	})

	m.Line(5)

	// Table Header
	header := []string{"Term", "Course Code", "Course Name", "Credits", "Grade", "Status"}
	m.Row(10, func() {
		m.Col(12, func() {
			m.TableList(header, [][]string{}, props.TableList{
				HeaderProp: props.TableListContent{
					Size:      9,
					GridSizes: []uint{2, 2, 4, 1, 1, 2},
				},
				ContentProp: props.TableListContent{
					Size:      8,
					GridSizes: []uint{2, 2, 4, 1, 1, 2},
				},
				Align:                consts.Center,
				// AlternatingColor:     &color.Color{Red: 240, Green: 240, Blue: 240},
				// HeaderBackgroundColor: &color.Color{Red: 200, Green: 200, Blue: 200},
			})
		})
	})

	// Table Content
	var tableData [][]string
	for _, e := range enrollments {
		term := "N/A"
		if e.Section != nil && e.Section.AcademicTerm != nil {
			term = e.Section.AcademicTerm.Name
		}
		code := "N/A"
		name := "N/A"
		credits := "0"
		if e.Section != nil && e.Section.Course != nil {
			code = e.Section.Course.Code
			name = e.Section.Course.Name
			credits = fmt.Sprintf("%d", e.Section.Course.Credits)
		}
		grade := "N/A"
		if e.Grade != nil {
			grade = fmt.Sprintf("%.2f", *e.Grade)
		}

		tableData = append(tableData, []string{
			term, code, name, credits, grade, e.Status,
		})
	}

	for _, row := range tableData {
		m.Row(8, func() {
			m.Col(12, func() {
				m.TableList([]string{}, [][]string{row}, props.TableList{
					ContentProp: props.TableListContent{
						Size:      8,
						GridSizes: []uint{2, 2, 4, 1, 1, 2},
					},
					Align: consts.Center,
				})
			})
		})
	}

	buf, err := m.Output()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
