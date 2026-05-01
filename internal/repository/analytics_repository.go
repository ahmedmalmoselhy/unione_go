package repository

import (
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetOverviewStats(facultyID, departmentID *uint) (models.OverviewStats, error)
	GetEnrollmentStatusStats(facultyID, departmentID *uint) (map[string]int64, error)
	GetGradeDistributionStats(facultyID, departmentID *uint) (map[string]int64, error)
	GetGPADistributionStats(facultyID, departmentID *uint) (map[string]int64, error)
	GetSectionEnrollmentRates(facultyID, departmentID *uint) ([]models.SectionEnrollmentRate, error)
	GetEnrollmentTrends(months int) ([]models.EnrollmentTrend, error)
	GetCourseDemand() ([]models.CourseDemand, error)
	GetProfessorWorkload() ([]models.ProfessorWorkload, error)
	GetAttendanceStats() ([]models.AttendanceStatusStat, float64, error)
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetOverviewStats(facultyID, departmentID *uint) (models.OverviewStats, error) {
	var stats models.OverviewStats

	// Students count
	query := r.db.Model(&models.User{}).Where("role = ?", models.RoleStudent)
	if facultyID != nil {
		query = query.Where("faculty_id = ?", *facultyID)
	}
	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}
	query.Count(&stats.Students)

	// Professors count
	profQuery := r.db.Model(&models.User{}).Where("role = ?", models.RoleProfessor)
	if facultyID != nil {
		profQuery = profQuery.Where("faculty_id = ?", *facultyID)
	}
	if departmentID != nil {
		profQuery = profQuery.Where("department_id = ?", *departmentID)
	}
	profQuery.Count(&stats.Professors)

	// Courses count
	courseQuery := r.db.Model(&models.Course{})
	if facultyID != nil {
		courseQuery = courseQuery.Joins("JOIN departments ON courses.department_id = departments.id").
			Where("departments.faculty_id = ?", *facultyID)
	}
	if departmentID != nil {
		courseQuery = courseQuery.Where("department_id = ?", *departmentID)
	}
	courseQuery.Count(&stats.Courses)

	// Sections count
	sectionQuery := r.db.Model(&models.Section{})
	if facultyID != nil {
		sectionQuery = sectionQuery.Joins("JOIN courses ON sections.course_id = courses.id").
			Joins("JOIN departments ON courses.department_id = departments.id").
			Where("departments.faculty_id = ?", *facultyID)
	}
	if departmentID != nil {
		sectionQuery = sectionQuery.Joins("JOIN courses ON sections.course_id = courses.id").
			Where("courses.department_id = ?", *departmentID)
	}
	sectionQuery.Count(&stats.Sections)

	return stats, nil
}

func (r *analyticsRepository) GetEnrollmentStatusStats(facultyID, departmentID *uint) (map[string]int64, error) {
	results := make(map[string]int64)
	
	type result struct {
		Status string
		Count  int64
	}
	var res []result
	
	enrollQuery := r.db.Model(&models.Enrollment{}).Select("status, count(*) as count")
	if facultyID != nil || departmentID != nil {
		enrollQuery = enrollQuery.Joins("JOIN users ON enrollments.student_id = users.id")
		if facultyID != nil {
			enrollQuery = enrollQuery.Where("users.faculty_id = ?", *facultyID)
		}
		if departmentID != nil {
			enrollQuery = enrollQuery.Where("users.department_id = ?", *departmentID)
		}
	}
	
	err := enrollQuery.Group("status").Scan(&res).Error
	if err != nil {
		return nil, err
	}

	for _, item := range res {
		results[item.Status] = item.Count
	}
	return results, nil
}

func (r *analyticsRepository) GetGradeDistributionStats(facultyID, departmentID *uint) (map[string]int64, error) {
	results := map[string]int64{
		"A": 0, "B": 0, "C": 0, "D": 0, "F": 0,
	}
	
	query := r.db.Model(&models.Enrollment{}).
		Select("grade").
		Where("grade IS NOT NULL")

	if facultyID != nil || departmentID != nil {
		query = query.Joins("JOIN users ON enrollments.student_id = users.id")
		if facultyID != nil {
			query = query.Where("users.faculty_id = ?", *facultyID)
		}
		if departmentID != nil {
			query = query.Where("users.department_id = ?", *departmentID)
		}
	}

	var grades []float64
	err := query.Pluck("grade", &grades).Error
	if err != nil {
		return nil, err
	}

	for _, g := range grades {
		if g >= 90 {
			results["A"]++
		} else if g >= 80 {
			results["B"]++
		} else if g >= 70 {
			results["C"]++
		} else if g >= 60 {
			results["D"]++
		} else {
			results["F"]++
		}
	}

	return results, nil
}

func (r *analyticsRepository) GetGPADistributionStats(facultyID, departmentID *uint) (map[string]int64, error) {
	brackets := map[string]int64{
		"0.0-1.99": 0,
		"2.0-2.49": 0,
		"2.5-2.99": 0,
		"3.0-3.49": 0,
		"3.5-4.0":  0,
	}

	// In this Go port, GPA is not stored in the User table. 
	// To implement this properly, we'd need to calculate GPA for all students or store it.
	// For now, returning empty to avoid complex on-the-fly calculations for all students in one query.
	return brackets, nil
}

func (r *analyticsRepository) GetSectionEnrollmentRates(facultyID, departmentID *uint) ([]models.SectionEnrollmentRate, error) {
	var sections []struct {
		ID         uint
		Capacity   int
		CourseCode string
		CourseName string
		Filled     int
	}

	query := r.db.Table("sections").
		Select("sections.id, sections.capacity, courses.code as course_code, courses.name as course_name, (SELECT COUNT(*) FROM enrollments WHERE enrollments.section_id = sections.id AND enrollments.status IN ('registered', 'completed', 'enrolled')) as filled").
		Joins("JOIN courses ON sections.course_id = courses.id").
		Joins("JOIN academic_terms ON sections.academic_term_id = academic_terms.id").
		Where("academic_terms.is_active = ?", true).
		Where("sections.deleted_at IS NULL")

	if facultyID != nil || departmentID != nil {
		query = query.Joins("JOIN departments ON courses.department_id = departments.id")
		if facultyID != nil {
			query = query.Where("departments.faculty_id = ?", *facultyID)
		}
		if departmentID != nil {
			query = query.Where("departments.id = ?", *departmentID)
		}
	}

	err := query.Scan(&sections).Error
	if err != nil {
		return nil, err
	}

	var rates []models.SectionEnrollmentRate
	for _, s := range sections {
		pct := 0.0
		if s.Capacity > 0 {
			pct = (float64(s.Filled) / float64(s.Capacity)) * 100
		}
		rates = append(rates, models.SectionEnrollmentRate{
			SectionID:  s.ID,
			CourseCode: s.CourseCode,
			CourseName: s.CourseName,
			Capacity:   s.Capacity,
			Filled:     s.Filled,
			FillPct:    pct,
		})
	}

	return rates, nil
}

func (r *analyticsRepository) GetEnrollmentTrends(months int) ([]models.EnrollmentTrend, error) {
	var trends []models.EnrollmentTrend
	
	// Check driver name to handle SQLite vs Postgres
	dialect := r.db.Dialector.Name()

	if dialect == "sqlite" {
		// Simplified for SQLite in tests
		err := r.db.Raw(`
			SELECT 
				strftime('%Y-%m-01 00:00:00', created_at) as month, 
				COUNT(*) as total_enrollments,
				COUNT(CASE WHEN status IN ('registered', 'enrolled') THEN 1 END) as active,
				COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
				COUNT(CASE WHEN status = 'dropped' THEN 1 END) as dropped
			FROM enrollments
			GROUP BY month
			ORDER BY month
		`).Scan(&trends).Error
		return trends, err
	}

	// Postgres syntax
	err := r.db.Raw(`
		SELECT 
			DATE_TRUNC('month', created_at) as month, 
			COUNT(*) as total_enrollments,
			COUNT(CASE WHEN status IN ('registered', 'enrolled') THEN 1 END) as active,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'dropped' THEN 1 END) as dropped
		FROM enrollments
		WHERE created_at >= NOW() - INTERVAL '1 month' * ?
		GROUP BY month
		ORDER BY month
	`, months).Scan(&trends).Error

	if err != nil {
		return nil, err
	}

	for i := range trends {
		if trends[i].TotalEnrollments > 0 {
			trends[i].DropRate = (float64(trends[i].Dropped) / float64(trends[i].TotalEnrollments)) * 100
		}
	}

	return trends, nil
}

func (r *analyticsRepository) GetCourseDemand() ([]models.CourseDemand, error) {
	var demand []models.CourseDemand

	err := r.db.Raw(`
		SELECT 
			courses.id as course_id,
			courses.code as code,
			courses.name as name,
			COUNT(enrollments.id) as current_enrollments,
			SUM(sections.capacity) as capacity,
			ROUND(COUNT(enrollments.id) * 100.0 / NULLIF(SUM(sections.capacity), 0), 2) as fill_rate
		FROM enrollments
		JOIN sections ON enrollments.section_id = sections.id
		JOIN courses ON sections.course_id = courses.id
		WHERE enrollments.status IN ('registered', 'enrolled')
		GROUP BY courses.id, courses.code, courses.name
		ORDER BY current_enrollments DESC
		LIMIT 20
	`).Scan(&demand).Error

	if err != nil {
		return nil, err
	}

	for i := range demand {
		if demand[i].FillRate > 90 {
			demand[i].DemandLevel = "high"
		} else if demand[i].FillRate > 70 {
			demand[i].DemandLevel = "medium"
		} else {
			demand[i].DemandLevel = "low"
		}
	}

	return demand, nil
}

func (r *analyticsRepository) GetProfessorWorkload() ([]models.ProfessorWorkload, error) {
	var workload []models.ProfessorWorkload

	err := r.db.Raw(`
		SELECT 
			users.id as professor_id,
			users.first_name || ' ' || users.last_name as name,
			'Professor' as rank,
			COUNT(sections.id) as total_sections,
			COUNT(CASE WHEN sections.deleted_at IS NULL THEN 1 END) as active_sections
		FROM users
		LEFT JOIN sections ON users.id = sections.professor_id
		WHERE users.role = 'professor'
		GROUP BY users.id, users.first_name, users.last_name
		ORDER BY active_sections DESC
	`).Scan(&workload).Error

	if err != nil {
		return nil, err
	}

	for i := range workload {
		if workload[i].ActiveSections > 5 {
			workload[i].WorkloadLevel = "heavy"
		} else if workload[i].ActiveSections > 3 {
			workload[i].WorkloadLevel = "moderate"
		} else {
			workload[i].WorkloadLevel = "light"
		}
	}

	return workload, nil
}

func (r *analyticsRepository) GetAttendanceStats() ([]models.AttendanceStatusStat, float64, error) {
	var stats []models.AttendanceStatusStat
	var total int64
	r.db.Model(&models.Attendance{}).Count(&total)

	if total == 0 {
		return stats, 0, nil
	}

	type res struct {
		IsPresent bool
		Count     int
	}
	var results []res
	err := r.db.Model(&models.Attendance{}).Select("is_present, count(*) as count").Group("is_present").Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	var presentCount int64
	for _, r := range results {
		status := "absent"
		if r.IsPresent {
			status = "present"
			presentCount = int64(r.Count)
		}
		stats = append(stats, models.AttendanceStatusStat{
			Status:     status,
			Count:      r.Count,
			Percentage: (float64(r.Count) / float64(total)) * 100,
		})
	}

	rate := (float64(presentCount) / float64(total)) * 100

	return stats, rate, nil
}
