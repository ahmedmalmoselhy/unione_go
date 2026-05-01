package models

import "time"

type DashboardStats struct {
	Overview          OverviewStats          `json:"overview"`
	EnrollmentStatus  map[string]int64       `json:"enrollment_status"`
	GradeDistribution map[string]int64       `json:"grade_distribution"`
	GPADistribution   map[string]int64       `json:"gpa_distribution"`
	EnrollmentRates   []SectionEnrollmentRate `json:"enrollment_rates"`
}

type OverviewStats struct {
	Students   int64 `json:"students"`
	Professors int64 `json:"professors"`
	Courses    int64 `json:"courses"`
	Sections   int64 `json:"sections"`
}

type SectionEnrollmentRate struct {
	SectionID  uint    `json:"section_id"`
	CourseCode string  `json:"course_code"`
	CourseName string  `json:"course_name"`
	Capacity   int     `json:"capacity"`
	Filled     int     `json:"filled"`
	FillPct    float64 `json:"fill_pct"`
}

type EnrollmentTrend struct {
	Month            time.Time `json:"month"`
	TotalEnrollments int       `json:"total_enrollments"`
	Active           int       `json:"active"`
	Completed        int       `json:"completed"`
	Dropped          int       `json:"dropped"`
	DropRate         float64   `json:"drop_rate"`
}

type EnrollmentTrendsResponse struct {
	Period string            `json:"period"`
	Data   []EnrollmentTrend `json:"data"`
}

type CourseDemand struct {
	CourseID           uint    `json:"course_id"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	CurrentEnrollments int     `json:"current_enrollments"`
	Capacity           int     `json:"capacity"`
	FillRate           float64 `json:"fill_rate"`
	DemandLevel        string  `json:"demand_level"`
}

type CourseDemandResponse struct {
	TopCourses []CourseDemand `json:"top_courses"`
}

type ProfessorWorkload struct {
	ProfessorID    uint   `json:"professor_id"`
	Name           string `json:"name"`
	Rank           string `json:"rank"`
	TotalSections  int    `json:"total_sections"`
	ActiveSections int    `json:"active_sections"`
	HoursPerWeek   int    `json:"hours_per_week"`
	WorkloadLevel  string `json:"workload_level"`
}

type ProfessorWorkloadSummary struct {
	AverageSections float64 `json:"average_sections"`
	AverageHours    float64 `json:"average_hours"`
	TotalProfessors int     `json:"total_professors"`
}

type ProfessorWorkloadResponse struct {
	Professors []ProfessorWorkload      `json:"professors"`
	Summary    ProfessorWorkloadSummary `json:"summary"`
}

type AttendanceStatusStat struct {
	Status     string  `json:"status"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type AttendanceAnalyticsResponse struct {
	ByStatus               []AttendanceStatusStat `json:"by_status"`
	OverallAttendanceRate float64                `json:"overall_attendance_rate"`
}

type StudentPerformancePrediction struct {
	Prediction   string  `json:"prediction"`
	Confidence   int     `json:"confidence"`
	CurrentGPA   float64 `json:"current_gpa"`
	AverageScore float64 `json:"average_score"`
	Trend        string  `json:"trend"`
	RiskLevel    string  `json:"risk_level"`
}
