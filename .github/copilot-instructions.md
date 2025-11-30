# Copilot Instructions for Simple Go CMS

## Architecture Overview

This is a **beginner-friendly** Go Clean Architecture template with **NO interfaces** - using direct struct dependencies for simplicity. Three-layer architecture:

**Handler → Usecase → Repository** (strict one-way dependency)

```
Handler (HTTP)  →  Usecase (Business Logic)  →  Repository (Database)
```

Key architectural decisions:
- **No interfaces**: Direct struct pointers (`*usecase.UserUsecase`, `*repository.UserRepository`) to reduce complexity for beginners
- **Fiber v2 framework**: Express.js-like API for HTTP handling
- **goqu query builder**: Lightweight, type-safe SQL building (NOT an ORM)
- **English comments**: All code comments are in English for better accessibility
- **Standardized responses**: Use `pkg/response` package for all HTTP responses

## Component Initialization Pattern

Follow this strict initialization order in `cmd/api/main.go`:

```go
// 1. Repository (bottom layer)
userRepo := repository.NewUserRepository(db)

// 2. Usecase (middle layer)
userUsecase := usecase.NewUserUsecase(userRepo)

// 3. Handler (top layer)
userHandler := handler.NewUserHandler(userUsecase)
```

Never skip layers - handlers must call usecases, usecases must call repositories.

## Adding New Features (CRUD Pattern)

When adding a new entity (e.g., "product"), create these files in order:

1. **Model** (`internal/model/product.go`): Define struct with JSON tags and request/response types
   ```go
   type Product struct {
       ID        int64     `db:"id" json:"id"`
       Name      string    `db:"name" json:"name"`
       CreatedAt time.Time `db:"created_at" json:"created_at"`
       UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
   }
   ```

2. **Repository** (`internal/repository/product_repository.go`): Database operations using goqu
   - Use `r.db.Dialect` for query building
   - Always use context: `r.db.SQL.ExecContext(ctx, query, args...)`
   - Return `fmt.Errorf("failed to X: %w", err)` for wrapped errors

3. **Usecase** (`internal/usecase/product_usecase.go`): Business logic and validation
   - Input validation happens here
   - String normalization (trim, lowercase emails)
   - Call repository methods, handle errors

4. **Handler** (`internal/handler/product_handler.go`): HTTP request/response
   - Parse request with `c.BodyParser(&req)` or `c.Params("id")`
   - Always use `response.*` helpers: `response.Success()`, `response.Created()`, `response.BadRequest()`, etc.
   - Never return raw JSON - always use response package

5. **Routes** in `cmd/api/main.go`: Register under `/api/v1/<resource>`
   ```go
   products := api.Group("/products")
   products.Get("/", productHandler.GetAll)
   products.Get("/:id", productHandler.GetByID)
   products.Post("/", productHandler.CreateProduct)
   ```

## Database Patterns

### Using goqu Query Builder

**Insert:**
```go
query, args, err := r.db.Dialect.
    Insert("table_name").
    Rows(goqu.Record{"field": value}).
    ToSQL()
result, err := r.db.SQL.ExecContext(ctx, query, args...)
```

**Select:**
```go
query, args, err := r.db.Dialect.
    Select("id", "name", "created_at").
    From("table_name").
    Where(goqu.Ex{"id": id}).
    ToSQL()
err = r.db.SQL.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.CreatedAt)
```

**Update:**
```go
query, args, err := r.db.Dialect.
    Update("table_name").
    Set(goqu.Record{"name": newName}).
    Where(goqu.Ex{"id": id}).
    ToSQL()
```

### Error Handling

- Use `sql.ErrNoRows` for "not found" cases:
  ```go
  if err == sql.ErrNoRows {
      return nil, fmt.Errorf("user not found")
  }
  ```
- Always wrap errors with context: `fmt.Errorf("failed to create user: %w", err)`

## Response Format

All responses follow this structure from `pkg/response`:

```json
{
  "success": true/false,
  "message": "User created successfully",
  "data": {...},
  "error": "error details (if failed)"
}
```

Use these helpers:
- `response.Success(c, data, "message")` - 200 OK
- `response.Created(c, data, "message")` - 201 Created
- `response.BadRequest(c, "message", err)` - 400 Bad Request
- `response.NotFound(c, "message")` - 404 Not Found
- `response.InternalServerError(c, "message", err)` - 500 Server Error

## Development Workflow

### Running the Application

**Local development:**
```bash
make run                  # Run without Docker
```

**Docker environment:**
```bash
make docker-up            # Start MySQL + API in containers
make docker-logs          # View logs
make docker-down          # Stop containers
make docker-rebuild       # Rebuild after code changes
```

### Database Migrations

Migrations are plain SQL files in `migrations/` directory:
- Automatically run on Docker startup (mounted to `/docker-entrypoint-initdb.d`)
- Manually run with: `make migrate` (requires `.env` file)

### Configuration

Environment variables loaded from `.env` (see `internal/config/config.go`):
- `SERVER_HOST`, `SERVER_PORT` - API server config
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Database connection
- `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS` - Connection pool settings

Default values exist for all configs, so `.env` is optional for local development.

## Code Conventions

- **Comments in English**: All inline comments and documentation are in English
- **No validation library**: Manual validation in usecases (see `validateCreateUser()` pattern)
- **No DTO layer**: Request/Response structs live in `internal/model/`
- **Context passing**: Always pass `c.Context()` from handler to usecase to repository
- **Struct naming**: `NewXxxHandler`, `NewXxxUsecase`, `NewXxxRepository` for constructors

## Common Pitfalls

❌ Don't create interfaces for repositories/usecases - this template intentionally avoids them
❌ Don't bypass layers (e.g., handler calling repository directly)
❌ Don't return raw Fiber errors - always use `pkg/response` helpers
❌ Don't forget `parseTime=true` in MySQL DSN (required for scanning timestamps)
❌ Don't use generic error messages - be specific about what failed

## Testing Endpoints

Health check:
```bash
curl http://localhost:8080/health
```

Sample API calls in README.md show complete CRUD examples for users endpoint.
