package services

import (
	"testing"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/repository"
	"gorm.io/gorm"
)

func migrateAll(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{}, &models.University{}, &models.Faculty{},
		&models.Department{}, &models.AcademicTerm{}, &models.Course{},
		&models.Section{}, &models.Enrollment{}, &models.Attendance{},
		&models.Announcement{}, &models.AnnouncementRead{},
		&models.Exam{}, &models.GroupProject{}, &models.GroupProjectMember{},
		&models.Waitlist{}, &models.CourseRating{}, &models.Notification{},
		&models.AuditLog{}, &models.Webhook{}, &models.WebhookDelivery{},
		&models.SectionTeachingAssistant{}, &models.StudentDepartmentHistory{},
	)
}

func TestAcademicService_GPA(t *testing.T) {
	db := setupTestDB(t)
	migrateAll(db)

	repo := repository.NewAcademicRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := NewAcademicService(repo, userRepo, nil, nil, nil, nil)

	// Setup data
	dept := models.Department{Name: "CS"}
	db.Create(&dept)
	
	course1 := models.Course{Code: "CS101", Name: "Intro", Credits: 3, DepartmentID: dept.ID}
	db.Create(&course1)
	course2 := models.Course{Code: "CS102", Name: "Adv", Credits: 4, DepartmentID: dept.ID}
	db.Create(&course2)
	
	term := models.AcademicTerm{Name: "Term 1"}
	db.Create(&term)
	
	section1 := models.Section{CourseID: course1.ID, AcademicTermID: term.ID}
	db.Create(&section1)
	section2 := models.Section{CourseID: course2.ID, AcademicTermID: term.ID}
	db.Create(&section2)
	
	student := models.User{Email: "s@s.com", Role: models.RoleStudent}
	db.Create(&student)
	
	gradeA := 4.0
	gradeB := 3.0
	db.Create(&models.Enrollment{StudentID: student.ID, SectionID: section1.ID, Grade: &gradeA, Status: "completed"})
	db.Create(&models.Enrollment{StudentID: student.ID, SectionID: section2.ID, Grade: &gradeB, Status: "completed"})

	gpa, err := svc.CalculateGPA(student.ID)
	if err != nil {
		t.Fatalf("CalculateGPA failed: %v", err)
	}

	expected := 24.0 / 7.0
	if gpa < expected-0.01 || gpa > expected+0.01 {
		t.Errorf("Expected GPA around %f, got %f", expected, gpa)
	}
}

func TestAcademicService_Prerequisites(t *testing.T) {
	db := setupTestDB(t)
	migrateAll(db)

	repo := repository.NewAcademicRepository(db)
	userRepo := repository.NewUserRepository(db)
	// We need a mock or real WebhookService/NotificationService if they are called
	// but we can pass nil if they are not used in these specific methods
	// Actually EnrollStudent uses webhookSvc
	svc := NewAcademicService(repo, userRepo, NewNotificationService(nil, nil, nil, nil, nil), nil, nil, NewWebhookService(db))

	dept := models.Department{Name: "CS"}
	db.Create(&dept)
	
	course1 := models.Course{Code: "CS101", Name: "Intro", Credits: 3, DepartmentID: dept.ID}
	db.Create(&course1)
	
	course2 := models.Course{Code: "CS102", Name: "Adv", Credits: 3, DepartmentID: dept.ID}
	db.Create(&course2)
	
	// Add prerequisite
	db.Model(&course2).Association("Prerequisites").Append(&course1)
	
	term := models.AcademicTerm{Name: "Term 1", IsActive: true}
	db.Create(&term)
	
	section2 := models.Section{CourseID: course2.ID, AcademicTermID: term.ID, Capacity: 10}
	db.Create(&section2)
	
	student := models.User{Email: "s2@s.com", Role: models.RoleStudent, NationalID: "STU2"}
	db.Create(&student)

	// Attempt to enroll in course2 without course1
	_, err := svc.EnrollStudent(student.ID, section2.ID)
	if err == nil {
		t.Error("Expected error enrolling in course with missing prerequisite, got nil")
	}

	// Complete course1
	section1 := models.Section{CourseID: course1.ID, AcademicTermID: term.ID}
	db.Create(&section1)
	grade := 80.0 // >= 60
	db.Create(&models.Enrollment{StudentID: student.ID, SectionID: section1.ID, Grade: &grade, Status: "completed"})

	// Now attempt again
	_, err = svc.EnrollStudent(student.ID, section2.ID)
	if err != nil {
		t.Errorf("Expected success enrolling after completing prerequisite, got error: %v", err)
	}
}
