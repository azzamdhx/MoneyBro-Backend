# MoneyBro Backend

Backend API untuk aplikasi manajemen keuangan pribadi MoneyBro.

## Tech Stack

- **Language**: Go 1.21+
- **API**: GraphQL (gqlgen)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Auth**: JWT
- **Email**: Resend

## Prerequisites

- Go 1.21 atau lebih baru
- PostgreSQL 15+
- Redis 7+
- Make (optional, untuk menjalankan commands)

## Getting Started

### 1. Clone & Install Dependencies

```bash
cd backend
go mod download
```

### 2. Setup Environment Variables

Copy file `.env.example` ke `.env` dan sesuaikan:

```bash
cp .env.example .env
```

```env
# Server
PORT=8080
ENV=development

# Database
DATABASE_URL=postgresql://postgres:password@localhost:5432/moneybro?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-super-secret-key-min-32-chars

# Resend (Email)
RESEND_API_KEY=re_xxxxxxxxxxxxxxxxxxxxxxxx

# Frontend URL (untuk CORS)
FRONTEND_URL=http://localhost:3000
```

### 3. Run Database Migrations

```bash
go run cmd/migrate/main.go up
```

### 4. Generate GraphQL Code

```bash
go generate ./...
```

### 5. Run Server

```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8080`

GraphQL Playground: `http://localhost:8080/playground`

## Project Structure

```
backend/
├── cmd/
│   ├── server/          # Main application entry point
│   └── migrate/         # Database migration tool
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database & Redis connections
│   ├── graph/           # GraphQL schema & resolvers
│   ├── middleware/      # HTTP middlewares (auth, cors, logging)
│   ├── models/          # Data models (GORM)
│   ├── repository/      # Database queries
│   ├── services/        # Business logic
│   └── utils/           # Helper functions
├── migrations/          # SQL migration files
├── .env.example
├── go.mod
├── go.sum
├── gqlgen.yml           # gqlgen configuration
├── Dockerfile
└── railway.toml
```

## Available Commands

```bash
# Run server
go run cmd/server/main.go

# Run migrations
go run cmd/migrate/main.go up
go run cmd/migrate/main.go down

# Generate GraphQL
go generate ./...

# Run tests
go test ./...

# Build for production
go build -o bin/server cmd/server/main.go
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/graphql` | POST | GraphQL API endpoint |
| `/playground` | GET | GraphQL Playground (dev only) |
| `/health` | GET | Health check |

## Deployment (Railway)

1. Connect repository ke Railway
2. Set environment variables di Railway dashboard
3. Railway akan auto-detect Go dan build menggunakan `Dockerfile`

## License

Private - Personal Use Only
