# CrazyLI – Agent-Orchestrated Google Login Demo

This repository hosts a minimal full-stack example powered by the OpenClaw orchestrator:

- **Backend**: Go + chi + Postgres (`backend/`)
- **Frontend**: React + Vite + TypeScript (`frontend/`)
- **Automation**: `.clawdbot/` scripts for spawning Codex/Claude/Gemini agents, monitoring, cron jobs, and CI blueprints.

## Getting Started

### 1. Provision Postgres

```bash
docker compose up -d postgres
```

_(A sample `docker-compose.yml` is coming soon. In the meantime you can run Postgres however you prefer and point `DATABASE_URL` at it.)_

Run the initial migration:

```bash
psql "$DATABASE_URL" -f backend/migrations/001_init.sql
```

### 2. Backend

```bash
cd backend
export DATABASE_URL=postgres://...
export GOOGLE_CLIENT_ID=...
export GOOGLE_CLIENT_SECRET=...
export GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
export FRONTEND_URL=http://localhost:5173
export JWT_SECRET=replace-me

go run ./cmd/server
```

### 3. Frontend

```bash
cd frontend
cp .env.example .env      # adjust VITE_API_BASE_URL if needed
pnpm install              # already run once, but safe to repeat
pnpm dev
```

Visit `http://localhost:5173` and click “使用 Google 登录”.

## Project Layout

```
backend/
  cmd/server          # HTTP server entrypoint
  internal/config     # env parsing
  internal/db         # Postgres helpers
  internal/auth       # Google OAuth + JWT
  migrations          # SQL schema
frontend/
  src/pages           # Home, success screen
  src/lib             # API helpers
.clawdbot/            # Agent automation scripts + templates
.github/workflows     # CI baseline
```

## Stage 5 Task Reference

The current feature – “Google 登录（Go 后端 + React 前端 + Postgres 存储）” – is the canonical example we’ll use to validate the Stage 5 workflow. Once Codex/Claude/Gemini CLIs are wired up, Zoe can reproduce/improve this feature via `.clawdbot` prompts.

## Authentication Modes
- **Google OAuth** (`/auth/google/login`) — existing flow.
- **Email + Password** (`/auth/email/register`, `/auth/email/login`) — stores bcrypt hashes, returns same session cookie.

Environment variables unchanged (JWT secret reused). To enable email login, simply run migrations `002` and `003`, restart backend, and use the new UI form on `/`.
