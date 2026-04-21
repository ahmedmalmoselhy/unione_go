# UniOne Platform - Go Implementation

Go backend implementation of UniOne using Gin, GORM, JWT-based auth, and layered handler/service/repository modules.

## Project Overview

The current Go port implements a partial UniOne backend with working coverage in these areas:
- authentication and JWT-protected identity lookup
- organization hierarchy CRUD for universities, faculties, and departments
- employee management under faculties
- Excel-based student import
- attendance recording and lookup
- GPA lookup
- grade import from uploaded files
- exam creation and listing
- announcement creation with notification service wiring

This repository is not yet at feature parity with the mature Laravel, Django, Node.js, or Rails ports.

## Technology Stack

- **Backend**: Go
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Auth**: JWT with role and scope middleware

## Current API Shape

The active router exposes:
- `/api/health`
- `/api/v1/auth/*`
- `/api/v1/organizations/*`
- `/api/v1/academic/*`
- `/api/v1/announcements`

Some additional handlers and models exist in the codebase, but the public route surface is still intentionally narrower than the mature ports.

## Documentation Index

- [README.md](README.md)
- [Enhancements.md](Enhancements.md)
- [API_ENDPOINTS.md](API_ENDPOINTS.md)
- [CURRENT_STATUS.md](CURRENT_STATUS.md)
- [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)

## Status

Partial backend implementation in progress.

Working foundations are in place, but major UniOne areas are still missing, including the broader student/professor/admin API surface, webhooks, analytics, exports, and test coverage.
