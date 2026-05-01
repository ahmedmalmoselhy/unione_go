package services

import (
	"fmt"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
)

type AnalyticsService interface {
	GetDashboardStats(facultyID, departmentID *uint) (*models.DashboardStats, error)
	GetEnrollmentTrends(months int) (*models.EnrollmentTrendsResponse, error)
	GetCourseDemand() (*models.CourseDemandResponse, error)
	GetProfessorWorkload() (*models.ProfessorWorkloadResponse, error)
	GetAttendanceAnalytics() (*models.AttendanceAnalyticsResponse, error)
	PredictStudentPerformance(studentID uint) (*models.StudentPerformancePrediction, error)
}

type analyticsService struct {
	repo        repository.AnalyticsRepository
	userRepo    repository.UserRepository
	academicSvc AcademicService
}

func NewAnalyticsService(repo repository.AnalyticsRepository, userRepo repository.UserRepository, academicSvc AcademicService) AnalyticsService {
	return &analyticsService{
		repo:        repo,
		userRepo:    userRepo,
		academicSvc: academicSvc,
	}
}

func (s *analyticsService) GetDashboardStats(facultyID, departmentID *uint) (*models.DashboardStats, error) {
	overview, err := s.repo.GetOverviewStats(facultyID, departmentID)
	if err != nil {
		return nil, err
	}

	enrollmentStatus, err := s.repo.GetEnrollmentStatusStats(facultyID, departmentID)
	if err != nil {
		return nil, err
	}

	gradeDist, err := s.repo.GetGradeDistributionStats(facultyID, departmentID)
	if err != nil {
		return nil, err
	}

	gpaDist, err := s.repo.GetGPADistributionStats(facultyID, departmentID)
	if err != nil {
		return nil, err
	}

	enrollmentRates, err := s.repo.GetSectionEnrollmentRates(facultyID, departmentID)
	if err != nil {
		return nil, err
	}

	return &models.DashboardStats{
		Overview:          overview,
		EnrollmentStatus:  enrollmentStatus,
		GradeDistribution: gradeDist,
		GPADistribution:   gpaDist,
		EnrollmentRates:   enrollmentRates,
	}, nil
}

func (s *analyticsService) GetEnrollmentTrends(months int) (*models.EnrollmentTrendsResponse, error) {
	trends, err := s.repo.GetEnrollmentTrends(months)
	if err != nil {
		return nil, err
	}

	return &models.EnrollmentTrendsResponse{
		Period: fmt.Sprintf("%d months", months),
		Data:   trends,
	}, nil
}

func (s *analyticsService) GetCourseDemand() (*models.CourseDemandResponse, error) {
	demand, err := s.repo.GetCourseDemand()
	if err != nil {
		return nil, err
	}

	return &models.CourseDemandResponse{
		TopCourses: demand,
	}, nil
}

func (s *analyticsService) GetProfessorWorkload() (*models.ProfessorWorkloadResponse, error) {
	workload, err := s.repo.GetProfessorWorkload()
	if err != nil {
		return nil, err
	}

	var totalSections, totalHours int
	for _, w := range workload {
		totalSections += w.ActiveSections
		totalHours += w.HoursPerWeek
	}

	avgSections := 0.0
	avgHours := 0.0
	if len(workload) > 0 {
		avgSections = float64(totalSections) / float64(len(workload))
		avgHours = float64(totalHours) / float64(len(workload))
	}

	return &models.ProfessorWorkloadResponse{
		Professors: workload,
		Summary: models.ProfessorWorkloadSummary{
			AverageSections: avgSections,
			AverageHours:    avgHours,
			TotalProfessors: len(workload),
		},
	}, nil
}

func (s *analyticsService) GetAttendanceAnalytics() (*models.AttendanceAnalyticsResponse, error) {
	stats, rate, err := s.repo.GetAttendanceStats()
	if err != nil {
		return nil, err
	}

	return &models.AttendanceAnalyticsResponse{
		ByStatus:               stats,
		OverallAttendanceRate: rate,
	}, nil
}

func (s *analyticsService) PredictStudentPerformance(studentID uint) (*models.StudentPerformancePrediction, error) {
	_, err := s.userRepo.FindByID(studentID)
	if err != nil {
		return nil, err
	}

	gpa, err := s.academicSvc.CalculateGPA(studentID)
	if err != nil {
		return nil, err
	}

	prediction := "C"
	if gpa >= 3.5 {
		prediction = "A"
	} else if gpa >= 3.0 {
		prediction = "B"
	} else if gpa < 2.0 {
		prediction = "F"
	}

	return &models.StudentPerformancePrediction{
		Prediction:   prediction,
		Confidence:   70,
		CurrentGPA:   gpa,
		AverageScore: gpa * 25, // Mock conversion
		Trend:        "stable",
		RiskLevel:    "low",
	}, nil
}
