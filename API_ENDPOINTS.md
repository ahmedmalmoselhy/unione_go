# UniOne Go - API Endpoints

This document tracks the endpoints currently exposed by `internal/routes/routes.go`.

## Health

- `GET /api/health`

## Authentication

- `POST /api/v1/auth/login`
- `POST /api/v1/auth/register`
- `GET /api/v1/auth/me`

## Organization

All organization routes are mounted under `/api/v1/organizations` and require authentication unless noted otherwise.

### Universities

- `GET /universities`
- `POST /universities`
- `PUT /universities/:id`
- `DELETE /universities/:id`

### Faculties

- `GET /universities/:university_id/faculties`
- `POST /universities/:university_id/faculties`
- `PUT /faculties/:id`
- `DELETE /faculties/:id`

### Departments

- `GET /faculties/:faculty_id/departments`
- `POST /faculties/:faculty_id/departments`
- `PUT /departments/:id`
- `DELETE /departments/:id`

### Employees

- `GET /faculties/:faculty_id/employees`
- `POST /faculties/:faculty_id/employees`
- `PUT /employees/:id`
- `DELETE /employees/:id`

### Imports

- `POST /faculties/:faculty_id/students/import`

## Academic

Academic routes are mounted under `/api/v1/academic`.

### Terms

- `GET /terms`
- `POST /terms`
- `GET /terms/:id`
- `PUT /terms/:id`
- `DELETE /terms/:id`

### Courses

- `GET /courses?department_id=:department_id`
- `POST /courses`
- `GET /courses/:id`
- `PUT /courses/:id`
- `DELETE /courses/:id`

### Sections

- `GET /sections?course_id=:course_id`
- `GET /sections?academic_term_id=:academic_term_id`
- `GET /sections?professor_id=:professor_id`
- `POST /sections`
- `GET /sections/:id`
- `PUT /sections/:id`
- `DELETE /sections/:id`

### Enrollments

- `POST /sections/:section_id/enrollments`
- `GET /sections/:section_id/enrollments`
- `DELETE /sections/:section_id/enrollments/:student_id`
- `GET /students/:student_id/enrollments`

### Grades

- `PATCH /sections/:section_id/students/:student_id/grade`

### Attendance

- `POST /sections/:section_id/attendance`
- `GET /sections/:section_id/attendance`

### GPA

- `GET /students/:student_id/gpa`

### Grade Import

- `POST /sections/:section_id/grades/import`

### Exams

- `POST /sections/:section_id/exams`
- `GET /sections/:section_id/exams`

## Announcements

- `POST /api/v1/announcements/`

## Admin

Admin routes are currently mounted under `/api/v1/admin`.

### Students

- `GET /students`
- `POST /students`
- `GET /students/:id`
- `PUT /students/:id`
- `DELETE /students/:id`
- `POST /students/:id/transfer`
- `GET /students/:id/transfers`

### Professors

- `GET /professors`
- `POST /professors`
- `GET /professors/:id`
- `PUT /professors/:id`
- `DELETE /professors/:id`

### Employees

- `GET /employees`
- `POST /employees`
- `GET /employees/:id`
- `PUT /employees/:id`
- `DELETE /employees/:id`

## Student Portal

Student portal routes are mounted under `/api/v1/student` and require a `student` role.

- `GET /profile`
- `GET /enrollments`
- `GET /transcript`
- `GET /academic-history`
- `GET /schedule`
- `GET /schedule/export`
- `GET /attendance`
- `GET /waitlist`
- `POST /waitlist`
- `DELETE /waitlist/:id`
- `POST /courses/:course_id/rate`

## Professor Portal

Professor portal routes are mounted under `/api/v1/professor` and require a `professor` role.

- `GET /profile`
- `GET /sections`
- `GET /sections/:id/roster`
- `PATCH /sections/:id/grade`
- `POST /sections/:id/attendance`
- `POST /sections/:id/announcements`
- `GET /schedule`

## Notes

- Several handlers exist for terms, courses, sections, enrollments, and grade updates, but they are not all mounted in the current router.
- Course prerequisites are now modeled and enforced during enrollment when prerequisite grades are available.
- Student transfer history is persisted through `StudentDepartmentHistory`.
- This document reflects exposed routes, not every handler present in the codebase.
