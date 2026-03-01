# Feature: Email + Password Login

## Requirement
- Add traditional email/password auth alongside Google OAuth on the login page.
- Users should be able to register with email + password and log in later.
- Frontend must expose both flows; backend must issue the same session cookies.

## Plan
1. **Schema**: add `password_hash` + `provider` columns, unique email index.
2. **Backend**:
   - Store helpers for email users (create, fetch).
   - Routes: `POST /auth/email/register`, `POST /auth/email/login`.
   - Reuse JWT issuance & `/me` response.
3. **Frontend**:
   - Extend Home page with email/password form (register or login toggle).
   - API helpers calling new endpoints.
4. **Testing**: curl / Postman for backend, browser for UI, ensure `/me` returns user.
5. **Docs**: update README + stage report.

## Execution Log
_(fill during development)_

## Testing
_(fill after verifications)_

## Release Notes
_(fill at the end)_
## Execution Log
- 03:42 — Added migration `002_email_auth.sql` (password_hash, provider, unique email). Applied via psql.
- 03:50 — Updated backend models + routes (`password.go`, email auth handlers, gofmt, go mod tidy).
- 04:12 — Added migration `003_google_id_nullable.sql` (allow NULL google_id).
- 04:20 — Added password utilities + server routes (`password.go`, `/auth/email/*`, `/me` provider field).
- 04:25 — Updated frontend (API helpers, dual login UI, CSS redesign, React entry) and rebuilt dev server.
- 04:32 — Fixed nullable google_id + scan helpers, added logging.

## Testing
- `curl` register/login happy path (email: test@example.com) — verified 200 responses + session cookie + `/me` payload with provider=email.
- Database spot check via `psql` to confirm new columns populated + password hash stored.
- Frontend manual test: reload `http://127.0.0.1:4173`, fill email form, observe success feedback + user card.

## Release Notes (draft)
- Added email/password authentication endpoints with JWT session cookie issuance.
- Frontend now shows Google + Email login cards with responsive layout.
- New migrations (`002_email_auth.sql`, `003_google_id_nullable.sql`).

## Monitoring / Retro
- Add `/auth/email/*` metrics (success/failure counts) to future logging pipeline.
- Update runbooks to include email credential reset steps.
