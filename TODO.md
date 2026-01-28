# Backend TODO

## 1. Project Setup

- [x] `go.mod` - Initialize Go module
- [x] `go.sum` - Dependencies lock file
- [x] `.env.example` - Environment variables template
- [x] `.gitignore` - Git ignore rules
- [x] `Makefile` - Common commands

### Folder Structure

- [x] `cmd/server/main.go`
- [x] `cmd/migrate/main.go`
- [x] `internal/config/config.go`
- [x] `internal/database/postgres.go`
- [x] `internal/database/redis.go`

---

## 2. Database Migrations

- [x] `migrations/000001_create_users_table.up.sql`
- [x] `migrations/000001_create_users_table.down.sql`
- [x] `migrations/000002_create_categories_table.up.sql`
- [x] `migrations/000002_create_categories_table.down.sql`
- [x] `migrations/000003_create_expenses_table.up.sql`
- [x] `migrations/000003_create_expenses_table.down.sql`
- [x] `migrations/000004_create_expense_templates_table.up.sql`
- [x] `migrations/000004_create_expense_templates_table.down.sql`
- [x] `migrations/000005_create_installments_table.up.sql`
- [x] `migrations/000005_create_installments_table.down.sql`
- [x] `migrations/000006_create_installment_payments_table.up.sql`
- [x] `migrations/000006_create_installment_payments_table.down.sql`
- [x] `migrations/000007_create_debts_table.up.sql`
- [x] `migrations/000007_create_debts_table.down.sql`
- [x] `migrations/000008_create_debt_payments_table.up.sql`
- [x] `migrations/000008_create_debt_payments_table.down.sql`
- [x] `migrations/000009_create_notification_logs_table.up.sql`
- [x] `migrations/000009_create_notification_logs_table.down.sql`
- [x] `migrations/000010_create_indexes.up.sql`
- [x] `migrations/000010_create_indexes.down.sql`

---

## 3. Models (GORM)

- [x] `internal/models/user.go`
- [x] `internal/models/category.go`
- [x] `internal/models/expense.go`
- [x] `internal/models/expense_template.go`
- [x] `internal/models/installment.go`
- [x] `internal/models/installment_payment.go`
- [x] `internal/models/debt.go`
- [x] `internal/models/debt_payment.go`
- [x] `internal/models/notification_log.go`

---

## 4. GraphQL Schema

- [x] `gqlgen.yml` - gqlgen configuration
- [x] `internal/graph/schema/schema.graphqls` - Root schema
- [x] `internal/graph/schema/user.graphqls`
- [x] `internal/graph/schema/category.graphqls`
- [x] `internal/graph/schema/expense.graphqls`
- [x] `internal/graph/schema/installment.graphqls`
- [x] `internal/graph/schema/debt.graphqls`
- [x] `internal/graph/schema/dashboard.graphqls`
- [x] Generate GraphQL code (`go run github.com/99designs/gqlgen generate`)

---

## 5. Repository Layer

- [x] `internal/repository/repository.go` - Interface definitions
- [x] `internal/repository/user_repository.go`
- [x] `internal/repository/category_repository.go`
- [x] `internal/repository/expense_repository.go`
- [x] `internal/repository/expense_template_repository.go`
- [x] `internal/repository/installment_repository.go`
- [x] `internal/repository/installment_payment_repository.go`
- [x] `internal/repository/debt_repository.go`
- [x] `internal/repository/debt_payment_repository.go`
- [x] `internal/repository/notification_log_repository.go`

---

## 6. Service Layer

- [x] `internal/services/services.go` - Service container
- [x] `internal/services/auth_service.go`
- [x] `internal/services/user_service.go`
- [x] `internal/services/category_service.go`
- [x] `internal/services/expense_service.go`
- [x] `internal/services/expense_template_service.go`
- [x] `internal/services/installment_service.go`
- [x] `internal/services/debt_service.go`
- [x] `internal/services/dashboard_service.go`
- [x] `internal/services/email_service.go`
- [x] `internal/services/notification_service.go`

---

## 7. GraphQL Resolvers

- [x] `internal/graph/resolver.go` - Resolver struct
- [x] `internal/graph/schema.resolvers.go` - All query/mutation resolvers
- [x] `internal/graph/converters.go` - Model converters
- [x] `internal/graph/scalar.go` - Custom scalars
- [x] `internal/graph/generated.go` - Generated code
- [x] `internal/graph/model/models_gen.go` - Generated models

---

## 8. Middleware

- [x] `internal/middleware/auth.go` - JWT authentication
- [x] `internal/middleware/cors.go` - CORS configuration
- [x] `internal/middleware/ratelimit.go` - Rate limiting

---

## 9. Utilities

- [x] `internal/utils/jwt.go` - JWT helpers
- [x] `internal/utils/password.go` - bcrypt helpers
- [x] `internal/utils/validator.go` - Input validation
- [x] `internal/utils/response.go` - Response helpers

---

## 10. Cron Jobs

- [x] `internal/cron/cron.go` - Cron scheduler with notification job

---

## 11. Deployment

- [x] `Dockerfile`
- [x] `railway.toml`
- [x] `docker-compose.yml`

---

## 12. Testing

- [ ] `internal/services/auth_service_test.go`
- [ ] `internal/services/expense_service_test.go`
- [ ] `internal/services/installment_service_test.go`
- [ ] `internal/services/debt_service_test.go`
- [ ] `internal/repository/user_repository_test.go`

---

## Progress Summary

| Section | Total | Done | Remaining |
|---------|-------|------|-----------|
| Project Setup | 9 | 9 | 0 |
| Migrations | 20 | 20 | 0 |
| Models | 9 | 9 | 0 |
| GraphQL Schema | 9 | 9 | 0 |
| Repository | 10 | 10 | 0 |
| Service | 11 | 11 | 0 |
| Resolvers | 6 | 6 | 0 |
| Middleware | 3 | 3 | 0 |
| Utilities | 4 | 4 | 0 |
| Cron Jobs | 1 | 1 | 0 |
| Deployment | 3 | 3 | 0 |
| Testing | 5 | 0 | 5 |
| **TOTAL** | **90** | **85** | **5** |

---

## Build Status

✅ `go build ./...` - PASSED  
✅ `go vet ./...` - PASSED  
✅ `go mod tidy` - PASSED
