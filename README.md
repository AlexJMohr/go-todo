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

## Deployment

The app is deployed to AWS Lambda + [Neon](https://neon.tech) serverless Postgres, fronted by API Gateway and served at **[gotodo.amohr.net](https://gotodo.amohr.net)**.

| Component | Service |
|-----------|---------|
| Compute | AWS Lambda (arm64, `provided.al2023`) |
| Database | Neon serverless Postgres |
| Routing | AWS API Gateway HTTP API |
| DNS / CDN | Cloudflare |
| TLS | AWS ACM |
| Infrastructure | Terraform (`terraform/`) |
| CI/CD | GitHub Actions (push to `main` deploys) |

### Infrastructure provisioning

Requires [Terraform](https://developer.hashicorp.com/terraform/install), AWS credentials, a Neon API key, and a Cloudflare API token.

```sh
cd terraform
terraform init
terraform apply \
  -var="neon_api_key=<key>" \
  -var="neon_org_id=<org-id>" \
  -var="cloudflare_api_token=<token>"
```

### Code deployment

Pushing to `main` automatically builds and deploys via GitHub Actions. The workflow:

1. Builds a Linux arm64 binary
2. Zips it
3. Uploads to Lambda via `aws lambda update-function-code`

### GitHub Actions secrets

| Secret | Description |
|--------|-------------|
| `AWS_ACCESS_KEY_ID` | IAM user with `lambda:UpdateFunctionCode` permissions |
| `AWS_SECRET_ACCESS_KEY` | Corresponding secret key |
| `AWS_REGION` | `us-west-2` |
