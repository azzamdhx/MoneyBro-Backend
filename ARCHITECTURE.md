# MoneyBro Backend - Architecture

## Overview

MoneyBro Backend menggunakan arsitektur **Clean Architecture** dengan pemisahan layer yang jelas untuk maintainability dan testability.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP Layer                               │
│                    (Chi Router + Middleware)                     │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                       GraphQL Layer                              │
│                    (gqlgen Resolvers)                            │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐   │
│  │   Queries    │  │  Mutations   │  │    Subscriptions     │   │
│  └──────────────┘  └──────────────┘  └──────────────────────┘   │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Service Layer                              │
│                    (Business Logic)                              │
│                                                                  │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌──────────────┐  │
│  │ AuthSvc    │ │ ExpenseSvc │ │InstallSvc  │ │  DebtSvc     │  │
│  └────────────┘ └────────────┘ └────────────┘ └──────────────┘  │
│                                                                  │
│  ┌────────────┐ ┌────────────────────────────────────────────┐  │
│  │ EmailSvc   │ │           NotificationSvc                  │  │
│  └────────────┘ └────────────────────────────────────────────┘  │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Repository Layer                            │
│                    (Data Access)                                 │
│                                                                  │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌──────────────┐  │
│  │ UserRepo   │ │ExpenseRepo │ │InstallRepo │ │  DebtRepo    │  │
│  └────────────┘ └────────────┘ └────────────┘ └──────────────┘  │
└─────────────────────────────┬───────────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              ▼                               ▼
┌─────────────────────────┐     ┌─────────────────────────┐
│       PostgreSQL        │     │          Redis          │
│    (Primary Storage)    │     │    (Cache + Session)    │
└─────────────────────────┘     └─────────────────────────┘
```

## Layer Responsibilities

### 1. HTTP Layer (`cmd/server`, `internal/middleware`)

- HTTP server setup (Chi router)
- CORS configuration
- Request logging
- Panic recovery
- Rate limiting

### 2. GraphQL Layer (`internal/graph`)

- Schema definition (`.graphqls` files)
- Resolver implementations
- Input validation
- Error handling & formatting

### 3. Service Layer (`internal/services`)

- Business logic
- Data transformation
- Cross-cutting concerns
- Transaction management

### 4. Repository Layer (`internal/repository`)

- Database queries (GORM)
- Redis caching
- Data persistence

### 5. Models (`internal/models`)

- GORM models
- GraphQL types
- DTOs

---

## Database Schema

### Entity Relationship Diagram

```
┌─────────────┐       ┌─────────────────┐       ┌──────────────────────┐
│   users     │       │   categories    │       │      expenses        │
├─────────────┤       ├─────────────────┤       ├──────────────────────┤
│ id (PK)     │◄──┬───│ id (PK)         │◄──────│ id (PK)              │
│ email       │   │   │ user_id (FK)    │       │ user_id (FK)         │
│ password    │   │   │ name            │       │ category_id (FK)     │
│ name        │   │   │ created_at      │       │ item_name            │
│ created_at  │   │   └─────────────────┘       │ unit_price           │
│ updated_at  │   │                             │ quantity             │
└─────────────┘   │                             │ notes                │
      │           │                             │ expense_date         │
      │           │                             │ created_at           │
      │           │                             └──────────────────────┘
      │           │
      │           │   ┌──────────────────────┐  ┌──────────────────────┐
      │           │   │  expense_templates   │  │ installments         │
      │           │   ├──────────────────────┤  ├──────────────────────┤
      │           └───│ id (PK)              │  │ id (PK)              │
      │               │ user_id (FK)         │  │ user_id (FK)         │◄─┐
      ├───────────────│ category_id (FK)     │  │ name                 │  │
      │               │ item_name            │  │ actual_amount        │  │
      │               │ unit_price           │  │ loan_amount          │  │
      │               │ quantity             │  │ monthly_payment      │  │
      │               │ recurring_day        │  │ tenor                │  │
      │               │ created_at           │  │ start_date           │  │
      │               └──────────────────────┘  │ due_day              │  │
      │                                         │ status               │  │
      │                                         │ notes                │  │
      │                                         │ created_at           │  │
      │                                         └──────────────────────┘  │
      │                                                                   │
      │               ┌──────────────────────┐  ┌──────────────────────┐  │
      │               │       debts          │  │installment_payments  │  │
      │               ├──────────────────────┤  ├──────────────────────┤  │
      │               │ id (PK)              │  │ id (PK)              │  │
      └───────────────│ user_id (FK)         │  │ installment_id (FK)  │──┘
                      │ person_name          │  │ payment_number       │
                      │ actual_amount        │  │ amount               │
                      │ loan_amount          │  │ paid_at              │
                      │ payment_type         │  │ created_at           │
                      │ monthly_payment      │  └──────────────────────┘
                      │ tenor                │
                      │ due_date             │  ┌──────────────────────┐
                      │ status               │  │   debt_payments      │
                      │ notes                │  ├──────────────────────┤
                      │ created_at           │  │ id (PK)              │
                      └──────────────────────┘  │ debt_id (FK)         │
                              │                 │ payment_number       │
                              └────────────────►│ amount               │
                                                │ paid_at              │
                                                │ created_at           │
                                                └──────────────────────┘

┌──────────────────────┐
│  notification_logs   │
├──────────────────────┤
│ id (PK)              │
│ user_id (FK)         │
│ type                 │
│ reference_id         │
│ sent_at              │
│ email_subject        │
│ created_at           │
└──────────────────────┘
```

---

## Authentication Flow

```
┌──────────┐      ┌──────────┐      ┌──────────┐      ┌──────────┐
│  Client  │      │  Server  │      │  Redis   │      │ Postgres │
└────┬─────┘      └────┬─────┘      └────┬─────┘      └────┬─────┘
     │                 │                 │                 │
     │  Login Request  │                 │                 │
     │────────────────►│                 │                 │
     │                 │                 │                 │
     │                 │  Verify User    │                 │
     │                 │─────────────────────────────────►│
     │                 │                 │                 │
     │                 │◄─────────────────────────────────│
     │                 │                 │                 │
     │                 │  Store Session  │                 │
     │                 │────────────────►│                 │
     │                 │                 │                 │
     │  JWT Token      │                 │                 │
     │◄────────────────│                 │                 │
     │                 │                 │                 │
     │  API Request    │                 │                 │
     │  + JWT Header   │                 │                 │
     │────────────────►│                 │                 │
     │                 │                 │                 │
     │                 │  Validate Token │                 │
     │                 │────────────────►│                 │
     │                 │                 │                 │
     │                 │◄────────────────│                 │
     │                 │                 │                 │
     │  Response       │                 │                 │
     │◄────────────────│                 │                 │
     │                 │                 │                 │
```

---

## Email Notification System

### Cron Job Flow

```
┌───────────────────────────────────────────────────────────────────┐
│                     Daily Cron Job (08:00 WIB)                    │
└─────────────────────────────┬─────────────────────────────────────┘
                              │
                              ▼
┌───────────────────────────────────────────────────────────────────┐
│                    NotificationService.Run()                       │
│                                                                    │
│  1. Get current date                                              │
│  2. Calculate H-1, H-2, H-3 dates                                 │
│                                                                    │
└─────────────────────────────┬─────────────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              ▼                               ▼
┌─────────────────────────┐     ┌─────────────────────────┐
│  Check Installments     │     │     Check Debts         │
│  WHERE due_day matches  │     │  WHERE due_date matches │
│  H-1, H-2, or H-3       │     │  H-1, H-2, or H-3       │
└───────────┬─────────────┘     └───────────┬─────────────┘
            │                               │
            └───────────────┬───────────────┘
                            ▼
┌───────────────────────────────────────────────────────────────────┐
│                    For each reminder:                              │
│                                                                    │
│  1. Check notification_logs (avoid duplicate)                     │
│  2. If not sent today → Send email via Resend                     │
│  3. Log to notification_logs                                      │
│                                                                    │
└───────────────────────────────────────────────────────────────────┘
```

### Email Templates

| Type | Subject | Trigger |
|------|---------|---------|
| `INSTALLMENT_REMINDER` | "Reminder: Cicilan {name} jatuh tempo dalam {n} hari" | H-3, H-2, H-1 |
| `DEBT_REMINDER` | "Reminder: Hutang ke {person} jatuh tempo dalam {n} hari" | H-3, H-2, H-1 |

---

## Caching Strategy (Redis)

| Key Pattern | TTL | Description |
|-------------|-----|-------------|
| `session:{user_id}` | 24h | User session data |
| `dashboard:{user_id}` | 5m | Dashboard summary cache |
| `user:{user_id}` | 1h | User profile cache |

### Cache Invalidation

- **Dashboard cache**: Invalidate on any expense/installment/debt mutation
- **User cache**: Invalidate on profile update

---

## Error Handling

### GraphQL Error Format

```json
{
  "errors": [
    {
      "message": "Human readable message",
      "path": ["mutation", "createExpense"],
      "extensions": {
        "code": "VALIDATION_ERROR",
        "field": "unit_price"
      }
    }
  ]
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHENTICATED` | 401 | Missing or invalid JWT |
| `FORBIDDEN` | 403 | Not authorized to access resource |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Input validation failed |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Security Considerations

1. **Password Hashing**: bcrypt with cost 12
2. **JWT**: RS256 or HS256 with 24h expiry
3. **CORS**: Whitelist frontend domain only
4. **Rate Limiting**: 100 req/min per IP
5. **Input Validation**: All inputs validated at GraphQL layer
6. **SQL Injection**: Prevented by GORM parameterized queries

---

## Performance Optimizations

1. **Database Indexes**:
   - `users(email)` - Unique index
   - `expenses(user_id, expense_date)` - Composite index
   - `installments(user_id, status)` - Composite index
   - `debts(user_id, status)` - Composite index

2. **N+1 Prevention**: DataLoader for GraphQL relations

3. **Connection Pooling**: GORM connection pool (max 25 connections)

4. **Redis Caching**: Dashboard and frequently accessed data
