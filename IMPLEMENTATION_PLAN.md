# UniOne Go Implementation Plan

Last Updated: April 22, 2026
Goal: bring the Go backend to backend parity with the completed UniOne ports.

## Objective

The Go port already has a working foundation. This plan turns that foundation into a full backend implementation that is comparable to Laravel, Django, Node.js, and Rails in feature coverage, API shape, and operational readiness.

This is a backend parity plan, not a frontend plan.

## Reference Repositories

Use these repositories as the source of truth while implementing parity:

- `unione_backend`: canonical feature and behavior reference
- `unione_django`: clean backend domain partitioning
- `unione_node`: infrastructure, monitoring, privacy, and queue-oriented patterns
- `unione_rails`: service-oriented organization and enhancement layering

## Current Starting Point

The Go codebase already includes:

- Gin router setup
- GORM and DB bootstrap
- JWT authentication and role/scope middleware
- organization CRUD for universities, faculties, and departments
- employee CRUD under faculties
- student Excel import
- attendance record/list endpoints
- GPA lookup
- grade import
- exam create/list endpoints
- announcement creation with notification service wiring

The major gap is not basic scaffolding anymore. The gap is missing breadth across the student, professor, shared, and admin domains.

## Parity Definition

The Go backend is considered on par when it provides:

- the core feature baseline from `UniOne/Features.md`
- stable `/api/v1` coverage for auth, student, professor, shared, and admin domains
- imports and exports
- notifications and email-driven flows
- exam schedule, teaching assistant, and group project management
- transcript and academic history flows
- waitlist, ratings, and prerequisite-aware enrollment
- audit logging and webhooks
- meaningful automated test coverage
- documentation that matches the implemented code

## Delivery Principles

1. Keep Laravel behavior as the primary parity target when there is ambiguity.
2. Do not expand routes faster than the service and repository layers can support cleanly.
3. Prefer one stable API shape under `/api/v1` rather than duplicating immature endpoints.
4. Build feature groups end to end: model, repository, service, handler, route, validation, tests, docs.
5. Add operational features only after the corresponding business flows are stable.

## Phased Plan

### Phase 0: Foundation Stabilization

Purpose: make the current codebase safe to build on.

Deliverables:

- standardize all public routes under `/api/v1`
- normalize request validation and error responses
- standardize auth context extraction from middleware
- verify role and scope middleware behavior
- add structured logging conventions
- define migration workflow and schema ownership
- document route conventions and response format

Exit criteria:

- one consistent API structure
- no handler-specific ad hoc response format
- middleware contract is stable enough for broader route expansion

### Phase 1: Core Academic CRUD

Purpose: complete the base academic model before portal features expand.

Deliverables:

- expose and complete academic term endpoints
- expose and complete course endpoints
- expose and complete section endpoints
- expose and complete enrollment endpoints
- expose grade update endpoints
- implement professor assignment to sections
- implement schedule persistence and retrieval
- implement capacity checks
- implement prerequisite validation

Exit criteria:

- terms, courses, sections, enrollments, and grading have a stable public API
- academic models support later student/professor/admin features without major redesign

### Phase 2: People Management

Purpose: complete the core identity-linked entities used across the product.

Deliverables:

- student CRUD
- professor CRUD
- extend employee management where needed for parity
- student department transfer and transfer history
- organization-aware filtering for all people endpoints

Exit criteria:

- admins can manage students, professors, and employees through stable APIs
- people records are correctly linked to organization structure and users

### Phase 3: Student Domain

Purpose: deliver the student backend surface expected in the mature ports.

Deliverables:

- student profile endpoint
- enrollments list/create/drop
- transcript JSON
- academic history
- schedule endpoint
- ICS schedule export
- attendance summary
- waitlist join/list/leave
- waitlist auto-promotion when seats open
- course rating flows

Exit criteria:

- the student route group covers the major self-service academic flows
- enrollment rules behave consistently with the mature ports

### Phase 4: Professor Domain

Purpose: deliver the professor-facing academic operations.

Deliverables:

- professor profile endpoint
- professor section list
- section student roster
- grade submission and update
- attendance session lifecycle
- section announcement CRUD
- professor schedule

Exit criteria:

- the professor route group supports grading, attendance, and section communication

### Phase 5: Admin Academic Management

Purpose: close the biggest parity gap with the completed ports.

Deliverables:

- admin CRUD for faculties, departments, courses, sections, academic terms
- teaching assistant assignment endpoints
- exam schedule show/create/update/publish
- group project CRUD and membership management
- import template endpoints
- admin dashboard summary endpoints if required for parity

Exit criteria:

- admin APIs can manage the academic structure and operational workflows without relying on unfinished ad hoc routes

### Phase 6: Communication, Notifications, and Delivery

Purpose: support the user-facing workflows triggered by announcements, grades, and exam publication.

Deliverables:

- shared announcements list and read-tracking
- notifications inbox, mark-read, mark-all-read, delete
- email dispatch integration for announcements, exams, and grades
- background job strategy for async delivery
- section-scoped and faculty/global announcement targeting

Exit criteria:

- communication flows support both in-app and async delivery behavior

### Phase 7: Documents, Imports, and Exports

Purpose: match the operational capabilities of the mature ports.

Deliverables:

- refine student Excel import
- refine professor grade import
- export endpoints for students, professors, employees, enrollments, and grades
- transcript PDF generation
- reusable import/export services and templates

Exit criteria:

- admins and professors can import/export key academic data through stable APIs

### Phase 8: Audit, Webhooks, and Governance

Purpose: make the backend operationally comparable to the completed ports.

Deliverables:

- audit logging middleware or service
- audit log query endpoints
- webhook registration and ownership checks
- signed webhook delivery
- delivery history and retry behavior
- auto-disable logic for repeated failures

Exit criteria:

- important mutations are auditable
- external integrations can subscribe to supported events safely

### Phase 9: Production Hardening

Purpose: move from feature-complete to deployment-ready.

Deliverables:

- rate limiting
- health and monitoring endpoints
- config validation
- privacy/GDPR endpoints if required for parity target
- performance profiling and key query optimization
- Swagger/OpenAPI generation
- CI pipeline

Exit criteria:

- operational controls exist for running the service outside local development

### Phase 10: Testing and Validation

Purpose: prove parity instead of assuming it.

Deliverables:

- unit tests for services and repositories
- handler tests for auth, student, professor, shared, and admin domains
- integration tests for enrollment, grading, exams, imports, waitlist, announcements, and webhooks
- parity checklist against Laravel routes and feature behaviors
- docs refresh after implementation

Exit criteria:

- the Go port has meaningful automated confidence and a documented parity checklist

## Recommended Execution Order

Follow this order to minimize rework:

1. Phase 0: Foundation Stabilization
2. Phase 1: Core Academic CRUD
3. Phase 2: People Management
4. Phase 3: Student Domain
5. Phase 4: Professor Domain
6. Phase 5: Admin Academic Management
7. Phase 6: Communication, Notifications, and Delivery
8. Phase 7: Documents, Imports, and Exports
9. Phase 8: Audit, Webhooks, and Governance
10. Phase 9: Production Hardening
11. Phase 10: Testing and Validation

## Immediate Next Step

Start with Phase 0 and Phase 1 together:

- stabilize middleware and API conventions
- expose the existing unmounted academic handlers
- complete terms, courses, sections, enrollments, and grade updates

That work unlocks almost every later phase.

## Working Checklist

### Phase 0

- [x] Standardize `/api/v1` route layout
- [x] Normalize error response shape
- [x] Normalize success response shape where needed
- [x] Verify JWT middleware context contract
- [x] Verify role/scope middleware contract
- [x] Add structured logging baseline
- [x] Document API conventions

### Phase 1

- [x] Expose term endpoints
- [x] Expose course endpoints
- [x] Expose section endpoints
- [x] Expose enrollment endpoints
- [x] Expose grade update endpoint
- [x] Add professor assignment support
- [x] Add prerequisite enforcement
- [x] Add capacity enforcement
- [x] Add schedule support

### Phase 2

- [x] Add student CRUD
- [x] Add professor CRUD
- [x] Complete employee parity gaps
- [ ] Add student transfer/history

### Phase 3

- [ ] Student profile
- [ ] Student enrollments
- [ ] Student transcript
- [ ] Student academic history
- [ ] Student schedule
- [ ] Student ICS export
- [ ] Student attendance summary
- [ ] Student waitlist
- [ ] Student ratings

### Phase 4

- [ ] Professor profile
- [ ] Professor sections
- [ ] Professor section students
- [ ] Professor grading
- [ ] Professor attendance workflow
- [ ] Professor section announcements
- [ ] Professor schedule

### Phase 5

- [ ] Admin academic terms
- [ ] Admin courses
- [ ] Admin sections
- [ ] Admin faculties/departments
- [ ] Teaching assistants
- [ ] Exam schedule publish flow
- [ ] Group projects
- [ ] Admin imports/templates

### Phase 6

- [ ] Shared announcements list/read
- [ ] Notifications inbox
- [ ] Mark read / mark all read
- [ ] Delete notification
- [ ] Email dispatch integration
- [ ] Background delivery workflow

### Phase 7

- [ ] Students import hardening
- [ ] Grades import hardening
- [ ] Export endpoints
- [ ] Transcript PDF
- [ ] Import/export templates

### Phase 8

- [ ] Audit log model/service
- [ ] Audit log endpoints
- [ ] Webhook registration
- [ ] Webhook delivery history
- [ ] Signed webhook delivery
- [ ] Retry and disable policy

### Phase 9

- [ ] Rate limiting
- [ ] Health/monitoring improvements
- [ ] Config validation
- [ ] Privacy endpoints if required
- [ ] Swagger/OpenAPI
- [ ] CI pipeline

### Phase 10

- [ ] Unit tests
- [ ] Handler tests
- [ ] Integration tests
- [ ] Parity checklist review
- [ ] Documentation refresh

## Risks

- exposing routes before the service layer is clean will create duplication
- role and scope behavior can drift if implemented ad hoc per handler
- imports, notifications, and webhooks will become hard to maintain if they stay synchronous and handler-driven
- parity claims will drift again if docs are not updated at the end of each completed phase

## Progress Tracking Rule

When a phase is completed:

1. update `CURRENT_STATUS.md`
2. update `API_ENDPOINTS.md`
3. update `Enhancements.md`
4. update central `UniOne` docs if the public status meaningfully changed
