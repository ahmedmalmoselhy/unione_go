# UniOne Go Implementation Plan

## Goal

Deliver a high-performance Go backend with equivalent domain coverage for UniOne, utilizing Go's concurrency strengths.

## Recommended Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) (Web Framework)
- **ORM**: [GORM](https://gorm.io/) (PostgreSQL driver)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **Authentication**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- **Documentation**: [swaggo/swag](https://github.com/swaggo/swag) (Swagger)
- **Testing**: Built-in `testing` package + [testify](https://github.com/stretchr/testify)
- **Background Jobs**: Goroutines (internal) or [asynq](https://github.com/hibiken/asynq) + Redis (for scaled processing)

## Proposed Repository Structure

```text
unione_go/
  cmd/
    api/
      main.go         # Entry point
  internal/
    models/           # GORM structs
    handlers/         # HTTP controllers
    services/         # Business logic
    repository/       # Database access layer
    middleware/       # Auth, RBAC, Logging
    routes/           # Route definitions
    config/           # Environment and DB config
  pkg/
    utils/            # Helper functions (Excel, PDF, iCal)
  docs/               # Swagger documentation
```

## API Organization

Following the standard route structure:
- `/api/auth/*`
- `/api/student/*`
- `/api/professor/*`
- `/api/admin/*`
- `/api/announcements/*`
- `/api/notifications/*`

## Phased Plan

### Phase 1: Core Foundation & Identity
- Initialize Go module and project structure.
- Configure PostgreSQL connection with GORM.
- Implement User models and JWT authentication.
- Implement Role-Based Access Control (RBAC) middleware.
- Add health check and basic auth endpoints.

### Phase 2: Organization Hierarchy
- Implement models: University, Faculty, Department.
- Scoped authorization (University/Faculty/Department level).
- Employee management and organization CRUDs.

### Phase 3: Academic Catalog & Enrollments
- Course and Section management.
- Student profiles and enrollment logic.
- Academic terms and schedule models.

### Phase 4: Professor Portal & Grading
- Professor section assignments.
- Grading system and attendance tracking.
- GPA calculation service.

### Phase 5: Documents & Exports
- PDF transcript generation.
- iCal (.ics) schedule export.
- Excel/CSV import for students and grades.

### Phase 6: Communication & Notifications
- Announcement system (Global and Section-scoped).
- Notification service (internal + email dispatch).
- Background workers for email delivery.

### Phase 7: Audit & Integrations
- Audit logging middleware (tracking changes).
- Webhook system (registration and delivery with retries).

### Phase 8: Testing & Optimization
- Unit tests for services and repositories.
- Integration tests for all API endpoints.
- Performance profiling and SQL optimization.

## Success Criteria

- 100% parity with the Laravel/Node.js API specification.
- Stable RBAC with proper scoping.
- High-concurrency support for grade processing and imports.
- Comprehensive Swagger documentation.
- >80% code coverage.
