# UniOne Go Implementation Plan

This implementation plan details the strategy for building the Go port of the UniOne backend, comprehensively mapped to the core requirements found in `Features.md`.

## Open Decisions Required

> [!IMPORTANT]
> Please review the framework decisions. We are planning to use **Gin** as the web framework and **GORM** as the ORM, as per the suggestions in the project documentation. Please let me know if you would prefer **Echo** or **Ent**.

## Proposed Architecture and Structure

We will adopt the widely accepted standard Go project layout for backend services:
- `cmd/api/main.go`: Main execution entry point.
- `internal/`: Private code specific to this app.
    - `models/`: Database schema definitions (GORM structs).
    - `handlers/`: HTTP handlers translating HTTP requests/responses (Controllers).
    - `services/`: Core business logic bridging handlers and database.
    - `routes/`: Routing and endpoint configuration.
    - `middlewares/`: JWT authentication, RBAC, and error handling.
- `pkg/`: Utility and helper functions (e.g., Excel parsers, email senders).

---

## Phased Implementation (Mapped to Features)

### Phase 1: Core Foundation & Identity (Features 3, 4, 10)
- Setup PostgreSQL connection using GORM.
- Scaffold JWT authentication and Role-Based Access Control (RBAC).
- Implement User models and their management endpoints (CRUD):
  - **Feature 3**: Create/Modify Students
  - **Feature 4**: Create/Modify Professors
  - **Feature 10**: Manage Employees

### Phase 2: Organization Hierarchy & Academic Catalog (Features 9, 17)
- Define the core institutional structure: University -> Faculty -> Department.
- **Feature 17**: Group students by University/Faculty/Department/Course.
- **Feature 9**: Manage courses, lectures, and their prerequisites.

### Phase 3: Assignments & Operations (Features 5, 6, 7, 8)
- Structure sections and assign schedules.
- **Feature 5**: Create/Modify lectures schedule (JSON sections).
- **Feature 6**: Assign professors to courses/sections.
- **Feature 7**: Assign teaching assistants to courses.
- **Feature 8**: Create/Manage group projects and members.

### Phase 4: Assessment & Grading (Features 13, 15)
- **Feature 13**: Publish exams schedule with conflict detection.
- **Feature 15**: Publish final grades (Professor submission).

### Phase 5: File Processing & Notifications (Features 1, 2, 11, 12, 14, 16)
- **Feature 1**: Employees upload student lists via Excel/CSV.
- **Feature 2**: Professors upload student grades via Excel/CSV.
- Notification System & Email Dispatchers:
  - **Feature 11**: Publish announcements (Global + Section-scoped).
  - **Feature 12**: Send announcements via email.
  - **Feature 14**: Send exams schedule via email upon publication.
  - **Feature 16**: Send final grades via email upon publication.

## Open Questions

> [!WARNING]
> 1. **Framework Confirmation:** Should we proceed with **Gin** + **GORM**, or would you prefer **Echo** + **Ent**?
> 2. **Background Jobs:** For phase 5 (Emailing and large Excel processing), should we leverage `goroutines` naturally, or set up a dedicated background task processor like `asynq` + Redis?

## Verification Plan

### Automated Tests
- We will write integration tests using Go's built-in testing features and `httptest` to cover all defined endpoints in `API_ENDPOINTS.md`.

### Manual Verification
- Boot the server locally and manually verify endpoints.
- Ensure Swagger generation or document integration to test the endpoints cleanly via tools like Postman/Insomnia.
