# Chirpy

Twitter-like REST API in Go with JWT auth and PostgreSQL.

## Stack

- Go 1.24
- PostgreSQL + sqlc
- JWT authentication
- Swagger docs

## Features

- User registration/login with refresh tokens
- Create, read, delete chirps (posts)
- User upgrades via Polka webhooks
- Protected routes with JWT middleware

## Setup

```bash
# Install dependencies
go mod download

# Configure environment
cp .env.example .env

# Run database
docker-compose up -d

# Start server
go run main.go
```

Server runs on `localhost:8080`.

## API Endpoints

**Auth**
- `POST /api/users` - Register
- `POST /api/login` - Login
- `POST /api/refresh` - Refresh token
- `POST /api/revoke` - Revoke token
- `PUT /api/users` - Update user (protected)

**Chirps**
- `GET /api/chirps` - List all chirps
- `GET /api/chirps/{id}` - Get chirp
- `POST /api/chirps` - Create chirp (protected)
- `DELETE /api/chirps/{id}` - Delete chirp (protected)

**Admin**
- `GET /admin/metrics` - View metrics
- `POST /admin/reset` - Reset database
- `GET /api/healthz` - Health check

**Webhooks**
- `POST /api/polka/webhooks` - User upgrade webhook

## Project Structure

```
.
├── internal/
│   ├── auth/          # JWT + password hashing
│   ├── config/        # App configuration
│   ├── database/      # sqlc generated code
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # Auth & metrics middleware
│   └── services/      # Business logic
├── sql/
│   ├── queries/       # SQL queries
│   └── schema/        # Database migrations
└── static/            # Static assets
```