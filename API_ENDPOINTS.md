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

## Notes

- Several handlers exist for terms, courses, sections, enrollments, and grade updates, but they are not all mounted in the current router.
- This document reflects exposed routes, not every handler present in the codebase.
