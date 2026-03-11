# Deferred Items - Phase 09

## Pre-existing Test Failures

1. **TestSecureHeaders_ProductionHSTS** (`backend/internal/adapters/web/middleware/secureheaders_test.go:56`)
   - HSTS header not set in production mode test
   - From Phase 09-04 security headers middleware
   - Not caused by 09-05 changes
