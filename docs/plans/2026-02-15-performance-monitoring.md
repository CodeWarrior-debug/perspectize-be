# Performance Monitoring Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Establish performance baselines and monitoring across the full stack — HTTP timing, DB query logging, GraphQL operation timing, Go benchmarks, and frontend Web Vitals.

**Architecture:** Lightweight observability layer using slog structured logging (no new infrastructure). GORM callbacks for DB timing, chi middleware for HTTP timing, gqlgen extensions for GraphQL timing, web-vitals npm package for frontend metrics.

**Tech Stack:** Go slog, chi middleware, GORM callbacks, gqlgen extensions, web-vitals (npm)

---

**Full plan details:** `.planning/phases/07.4-performance-monitoring/07.4-01-PLAN.md`

### Task 1: Request timing middleware (backend/pkg/middleware/timing.go)
### Task 2: GORM slow query logger + DB stats endpoint (backend/pkg/database/)
### Task 3: GraphQL operation timing extension (backend/pkg/graphql/timing.go)
### Task 4: Wire all backend monitoring into main.go
### Task 5: Go benchmark tests for service layer
### Task 6: Frontend Web Vitals baseline (frontend/src/lib/vitals.ts)
