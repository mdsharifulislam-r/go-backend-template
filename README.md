# Go Backend Template

Professional Go API template inspired by the NestJS modular architecture. Uses **PostgreSQL** with **GORM** (the easiest Go ORM for Postgres).

## Stack

- Go + [chi](https://github.com/go-chi/chi) router
- PostgreSQL + [GORM](https://gorm.io)
- JWT auth, bcrypt, request validation
- Structured JSON responses (Nest-style)
- SMTP email (OTP logged to console in development)
- Docker Compose for Postgres

## Project structure

```text
cmd/api/                 # application entrypoint
config/                  # env configuration
internal/
  database/              # GORM + Postgres connection & auto-migrate
  models/                # GORM models
  module/
    auth/                # login, OTP, password flows
    user/                # register, profile, listing
  middleware/            # auth, validate, logger, cors, upload
  email/                 # SMTP + templates
  helper/                # jwt, bcrypt, otp
  querybuilder/          # search / filter / paginate
  response/              # sendResponse helper
  errors/                # AppError + handler wrapper
  server/                # HTTP server wiring
uploads/                 # uploaded files
```

## Quick start

1. Start Postgres:

```bash
docker compose up -d
```

2. Copy env:

```bash
cp .env.example .env
```

3. Install deps & run:

```bash
go mod tidy
go run ./cmd/api
```

API base URL: `http://localhost:8080/api/v1`

## Main endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/user/` | No | Register user |
| GET | `/api/v1/user/profile` | Bearer | Get profile |
| PATCH | `/api/v1/user/profile` | Bearer | Update profile (multipart) |
| GET | `/api/v1/user/all` | No | List users (`page`, `limit`, `searchTerm`) |
| POST | `/api/v1/auth/login` | No | Login |
| POST | `/api/v1/auth/verify-otp` | No | Verify email OTP |
| POST | `/api/v1/auth/forgot-password` | No | Send reset OTP |
| POST | `/api/v1/auth/reset-password` | Reset token header | Reset password |
| POST | `/api/v1/auth/change-password` | Bearer | Change password |

### Register example

```bash
curl -X POST http://localhost:8080/api/v1/user/ \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"password123"}'
```

OTP is printed in the server console when email credentials are empty.

### Login example

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'
```

## Notes

- Super admin is seeded from `SUPER_ADMIN_EMAIL` / `SUPER_ADMIN_PASSWORD` on startup.
- GORM `AutoMigrate` creates/updates tables in development.
- Keep `JWT_SECRET` strong in production and set `DB_SSLMODE` appropriately.
