# Startup Idea Graveyard

> Bury your abandoned startup ideas with a cause of death and an epitaph — then inspect the autopsy report to see why your ideas keep dying.

## Architecture

```mermaid
flowchart LR
    User["Browser\n(React + Vite + Tailwind + Framer Motion)"]
    API["Go + Gin API\n(/api/v1, /api/docs, /metrics)"]
    DB[("PostgreSQL 16\n(pgx/v5 + GIN Full-Text)")]

    User -- "JSON + In-Memory JWT\n+ httpOnly Refresh Cookie" --> API
    API -- "Context-scoped Queries\n& Transactions" --> DB
```

- **Backend**: Go, Gin, PostgreSQL 16 (`pgx/v5`), `golang-migrate`, `argon2id` password hashing, 15m JWT access tokens + 7d rotating `httpOnly` refresh cookies, `log/slog` structured JSON logs, Prometheus `/metrics`, Swagger UI at `/api/docs/index.html`.
- **Frontend**: React 18, TypeScript, Vite, React Router (single-page scrolling layout with overlay routes), Tailwind CSS, Framer Motion, and a minimal native `fetch` client. Zero extra UI, state, or chart libraries.

---

## Quick Start (Docker Compose)

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
2. Start PostgreSQL, run migrations, and boot the API + Frontend:
   ```bash
   docker compose up --build
   ```
3. Seed the demo graveyard (25 realistic dead startup ideas):
   ```bash
   make seed
   ```
4. Open:
   - **Web App**: `http://localhost:5173`
   - **Swagger Docs**: `http://localhost:8080/api/docs/index.html`
   - **Health / Readiness / Metrics**: `http://localhost:8080/healthz`, `/readyz`, `/metrics`

---

## Local Development (Without Docker for App)

### Backend
```bash
cd backend
export DATABASE_URL="postgres://gravehard:gravehard@localhost:5432/gravehard?sslmode=disable"
go run ./cmd/api
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

---

## Testing & Linting

- **Run all backend unit & integration tests**:
  ```bash
  make test
  ```
  > Integration tests in `internal/http/*_integration_test.go` use `testcontainers-go` to spin up a real PostgreSQL 16 container when Docker is available.
- **Run backend linter**:
  ```bash
  make lint
  ```
- **Typecheck & build frontend**:
  ```bash
  cd frontend && npm run build
  ```

---

## Environment Variables

See [`.env.example`](.env.example) for the complete list:

| Variable | Description | Default |
|---|---|---|
| `APP_ENV` | `development` or `production` | `development` |
| `PORT` | API listen port | `8080` |
| `CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins | `http://localhost:5173` |
| `DATABASE_URL` | PostgreSQL connection string | *(required)* |
| `JWT_SECRET` | HMAC-SHA256 signing secret (32+ bytes in prod) | *(required in prod)* |
| `COOKIE_SECURE` | Set `true` behind HTTPS | `false` |
| `RATE_LIMIT_GENERAL_RPM` | Per-IP RPM for `/api/v1/*` | `120` |
| `RATE_LIMIT_AUTH_RPM` | Per-IP RPM for `/api/v1/auth/*` | `20` |
| `RATE_LIMIT_DEMO_RPM` | Per-IP RPM for `/api/v1/demo/*` | `60` |
| `VITE_API_URL` | Frontend API base URL | `http://localhost:8080/api/v1` |

---

## Key Design Decisions & Trade-offs

1. **`pgx/v5` + explicit repositories**: Provides compile-safe query control, native PostgreSQL `ENUM` / `death_cause[]` / stored `TSVECTOR` GIN full-text search, and clean transaction boundaries without ORM magic or Windows `sqlc` wasm issues.
2. **In-memory JWT + `httpOnly` Refresh Rotation**: Access tokens live strictly in React memory (immune to `localStorage` XSS exfiltration), while refresh tokens use `httpOnly`, `SameSite=Lax` cookies with SHA-256 DB hashing and reuse-detection family revocation.
3. **Integer Cents (`BIGINT`)**: All monetary values (`money_spent_cents`) are stored and transmitted as integer cents to eliminate floating-point rounding drift.
4. **Server-Side Pattern Engine (`n >= 3`)**: Failure observations are computed deterministically on the backend and suppressed when sample size `< 3` to avoid misleading insights on tiny graveyards.
5. **Single-Page Overlay Routing**: Keeping the entire product on one scrolling page while mounting `/login`, `/register`, `/bury`, and `/idea/:id` as modal overlays preserves scroll position and simplifies state sharing.

---

## Deployment

### Backend + PostgreSQL (Fly.io)
1. Provision Postgres & app in `backend/`:
   ```bash
   cd backend
   fly launch --no-deploy
   fly postgres create
   fly postgres attach <pg-app-name>
   fly secrets set JWT_SECRET="$(openssl rand -hex 32)" CORS_ALLOWED_ORIGINS="https://your-frontend.vercel.app"
   fly deploy
   ```
   `backend/fly.toml` is pre-configured with `/readyz` health checks and forced HTTPS.

### Frontend (Vercel / Netlify)
1. Set root directory to `frontend/`, build command `npm run build`, output directory `dist`.
2. Set `VITE_API_URL=https://gravehard-api.fly.dev/api/v1`.
3. `frontend/vercel.json` handles SPA rewrites for overlay routes.

---

## Security Checklist
- [x] `argon2id` password hashing with per-user random salt
- [x] No tokens in `localStorage` or `sessionStorage`
- [x] Refresh token rotation + family revocation on reuse
- [x] Per-IP token-bucket rate limiting (stricter on `/auth/*` and `/demo/*`)
- [x] Explicit CORS allowlist with credentials
- [x] 1 MiB request body limit + security headers (`X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`)
- [x] Scoped ownership checks (`404` instead of `403` for other users' idea IDs)
- [x] Non-root distroless Docker image
