---
created: 2026-03-11T00:13:10.660Z
title: Set up custom domain and Clerk production instance
area: deployment
files:
  - frontend/src/routes/+layout.svelte:8
  - backend/internal/config/security.go:45
---

## Problem

Clerk Development instance only allows `localhost` origins — the Sevalla-deployed frontend (`perspectize-fe-rf767.sevalla.page`) can't load Clerk JS/UI scripts because the dev instance blocks non-localhost domains. A production Clerk instance requires CNAME DNS records on a custom domain.

## Solution

1. Buy/configure a custom domain (e.g., `perspectize.com`)
2. Point the domain to Sevalla (both frontend and backend subdomains)
3. Create Clerk production instance with the custom domain
4. Add Clerk CNAME records to DNS:
   - `clerk.<domain>` → `frontend-api.clerk.services`
   - `accounts.<domain>` → `accounts.clerk.services`
   - Plus any additional CNAME records Clerk requires (check DNS Configuration page — 5 total)
5. Verify DNS in Clerk dashboard
6. Update Sevalla env vars with production Clerk keys:
   - Frontend: `PUBLIC_CLERK_PUBLISHABLE_KEY` → `pk_live_...`
   - Backend: `CLERK_SECRET_KEY` → `sk_live_...`
7. Update `CORS_ORIGINS` in backend to the custom domain
8. Set `APP_ENV=production` in backend to disable playground/introspection
