# Stage 5 Report – End-to-End Login Feature Demo

Date: 2026-02-28

## Objective
Implement a simple but complete “Google 登录” story (Go backend + React frontend + Postgres) to validate the orchestration-ready codebase before handing tasks to agents.

## Deliverables
1. **Backend (Go)**
   - Module initialised (`go.mod`).
   - `backend/cmd/server/main.go`: chi router, Google OAuth flow, JWT cookie issuance, `/me` endpoint.
   - `backend/internal/*`: config loader, Postgres store, OAuth helpers, JWT helpers.
   - `backend/migrations/001_init.sql`: `users` table.
   - `.env.example` for backend variables.

2. **Frontend (React + Vite)**
   - `frontend/` scaffolded via Vite + TypeScript + React Router.
   - Home page with Google login CTA + current user panel.
   - `/login/success` route to handle callback redirect.
   - `.env.example` with `VITE_API_BASE_URL`.

3. **Tooling**
   - `docker-compose.yml` (Postgres 16).
   - Updated root `README.md` with setup instructions & Stage 5 context.

## How to Test Manually
1. `docker compose up -d postgres`
2. `psql $DATABASE_URL -f backend/migrations/001_init.sql`
3. Export env vars (see `backend/.env.example`) and run `go run ./backend/cmd/server`
4. In another shell: `cd frontend && cp .env.example .env && pnpm dev`
5. Visit `http://localhost:5173`, click “使用 Google 登录”, complete OAuth. After redirect, `/me` returns the stored profile.

## Gaps / Next Steps
- Actual Google credentials must be provided by you.
- `.clawdbot` agents still need Codex/Claude CLI binaries before automating enhancements.
- Future tasks: persist sessions, add logout endpoint, enforce HTTPS in production.

This completes the Stage 5 dry-run requirement and seeds the repo with a concrete feature for future agent iterations.
