# go-todo

A simple todo app built with Go, Gin, GORM, HTMX, and PostgreSQL.

## Overview

| Component | Library |
|-----------|---------|
| Web framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) with [pgx](https://github.com/jackc/pgx) driver |
| Frontend interactivity | [htmx](https://htmx.org) |
| Database | PostgreSQL 17 |
| Live reload | [air](https://github.com/air-verse/air) |
| Env loading | [godotenv](https://github.com/joho/godotenv) |

## Development

### Prerequisites

- Go 1.25+
- Docker / Docker Compose
- [air](https://github.com/air-verse/air) (`go install github.com/air-verse/air@latest`)

### Setup

1. Copy the example env file and adjust if needed:

   ```sh
   cp .env.example .env
   ```

2. Start the database:

   ```sh
   docker compose up -d
   ```

3. Run the app with live reload:

   ```sh
   air
   ```

   The server starts on `http://localhost:8080` by default.

### Without live reload

```sh
go run .
```

### Environment variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/go_todo` |
