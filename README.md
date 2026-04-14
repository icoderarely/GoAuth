# GoAuth JWT

A small Go authentication API built to learn and explore core auth concepts:

- JWT-based authentication
- Password hashing with `bcrypt`
- Role-based authorization (`user`, `admin`)
- Clean separation of domain, service, handler, middleware, and repository layers

This project intentionally uses an **in-memory repository** so the focus stays on auth logic rather than database setup.

## Why an admin user is hard-coded in `main`

In `cmd/api/main.go`, an `admin` user is seeded on startup and promoted to `admin` role.

This is intentionally done for learning and quick experimentation:

- Lets you test admin-only routes immediately
- Demonstrates role assignment flow clearly
- Removes setup friction while exploring middleware and claims behavior

Important: this is **not a production pattern**. In real systems, use secure bootstrap/admin provisioning and secret management.

## What I learned from this project

- How JWT works: claims, signing, validation, token expiry
- Why passwords must be hashed, not stored in plain text
- How `bcrypt` helps protect user credentials
- How role checks can be enforced with middleware
- How to design auth-first application layers before introducing DB complexity

## Tech Stack

- Go
- `github.com/golang-jwt/jwt/v5`
- `golang.org/x/crypto/bcrypt`
- `github.com/google/uuid`
- `github.com/joho/godotenv`

## Project Structure

```text
cmd/api/main.go                 # app bootstrap
config/config.go                # environment config
internal/domain/                # entities + domain errors
internal/service/               # auth business logic
internal/repository/inmemory/   # in-memory user store
internal/handler/               # HTTP handlers
internal/middleware/            # auth + role middleware
internal/router/router.go       # route wiring
```

## Prerequisites

- Go 1.22+ (or compatible with your `go.mod`)

## Setup

1. Clone and enter the project.
2. Copy the example env file:

```bash
cp .env.example .env
```

3. Update `.env` values if needed.
4. Run:

```bash
go run ./cmd/api
```

Server starts on `http://localhost:8080` by default.

## Environment Variables

- `PORT` - API port (default: `8080`)
- `JWT_SECRET` - signing key for JWT tokens

## API Endpoints

### Public

- `POST /register`
- `POST /login`

### Protected

- `GET /me` (requires valid Bearer token)
- `GET /dashboard` (requires valid Bearer token)
- `GET /admin` (requires valid Bearer token with `admin` role)

## Quick Usage with cURL

Register a user:

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"alicepass"}'
```

Login:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"alicepass"}'
```

Login as seeded admin:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"pass"}'
```

Call protected endpoint:

```bash
curl http://localhost:8080/me \
  -H "Authorization: Bearer <TOKEN>"
```

Call admin endpoint:

```bash
curl http://localhost:8080/admin \
  -H "Authorization: Bearer <ADMIN_TOKEN>"
```

## Notes

- Data is stored in memory and resets when the server restarts.
- The hard-coded admin credentials are for local learning only.
- Improve next by adding persistent storage, user role management endpoints, refresh tokens, and stronger validation.
