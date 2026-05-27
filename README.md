# knowledgeKeeperApi

Personal knowledge-keeper REST API written in Go. Auth (register / login / refresh) backed by MariaDB, JWT access tokens with opaque refresh tokens, OpenAPI/Swagger UI bundled in the binary.

---

## Requirements

| Tool | Version | Why |
|---|---|---|
| Go | 1.25.3 | Build the API |
| Docker + Compose | recent | Run MariaDB (and optionally the API) locally |
| swag CLI | v1.16+ | Regenerate the OpenAPI spec from annotations (only needed when annotations change) |

Install `swag` once:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

Make sure `$(go env GOPATH)/bin` is on your `PATH` (typically `~/go/bin`):

```sh
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

---

## Configuration

The API reads everything from environment variables. A working sample lives in `.env`:

| Variable | Required | Default | Notes |
|---|---|---|---|
| `PORT` | no | `8080` | HTTP listen port |
| `DB_HOST` | yes | — | MariaDB host |
| `DB_PORT` | yes | — | MariaDB port |
| `DB_NAME` | yes | — | Database name |
| `DB_USER` | yes | — | App DB user |
| `DB_PASSWORD` | no | — | App DB password |
| `JWT_SECRET` | yes | — | HS256 signing secret. Use a long random string. |
| `JWT_ISSUER` | no | `knowledgeKeeperApi` | `iss` claim |
| `JWT_TTL` | no | `15m` | Access-token lifetime (Go duration: `15m`, `1h`, …) |
| `JWT_REFRESH_TTL` | no | `168h` | Refresh-token lifetime (7 days) |
| `MARIADB_ROOT_PASSWORD` | yes (compose) | — | Used by the MariaDB container only |

Generate a strong `JWT_SECRET`:

```sh
openssl rand -base64 48
```

---

## Run

### Option A — full stack with Docker Compose (recommended)

```sh
docker compose up --build
```

This starts MariaDB and the API. The first boot applies the auth schema automatically (`internal/auth/infra/sql/001_users.sql` is embedded in the binary).

API: <http://localhost:8080>
Swagger UI: <http://localhost:8080/swagger/index.html>

Stop everything (keeps DB volume):

```sh
docker compose down
```

Stop and wipe the DB volume:

```sh
docker compose down -v
```

### Option B — local Go, Dockerized DB

```sh
docker compose up -d db
DB_HOST=127.0.0.1 go run ./cmd/api
```

`DB_HOST=127.0.0.1` overrides the compose-internal hostname `db`.

### Option C — go run, no Docker

Bring your own MariaDB, then:

```sh
go run ./cmd/api
```

---

## Endpoints

| Method | Path | Auth | Purpose |
|---|---|---|---|
| `GET` | `/version` | none | Build version (CI-friendly) |
| `GET` | `/health` | none | Liveness probe |
| `POST` | `/auth/register` | none | Create a user |
| `POST` | `/auth/login` | none | Exchange credentials for tokens |
| `POST` | `/auth/refresh` | none | Rotate the token pair |
| `GET` | `/swagger/*` | none | Swagger UI + raw spec |

Quick smoke test:

```sh
curl -X POST localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"hunter22"}'

curl -X POST localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"hunter22"}'
```

---

## OpenAPI / Swagger

The spec is generated from godoc-style annotations on the handlers (`swaggo/swag`). Annotations are the source of truth; the `docs/` folder is generated and committed.

Regenerate after editing any annotation or adding a route:

```sh
./scripts/swag.sh
```

This produces `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`.

The Swagger UI is served by the binary itself at `/swagger/index.html` — no external files at runtime.

### Where annotations live

- General API metadata (title, version, security scheme): `cmd/api/main.go`
- Per-route annotations: `internal/auth/infra/http/handlers.go`, `internal/infra/http/server.go`
- DTOs referenced by `@Param` / `@Success` / `@Failure`: exported types in those same files

Bearer auth is declared globally (`@securityDefinitions.apikey BearerAuth`); apply it to a route by adding `// @Security BearerAuth` to the handler godoc once you start protecting endpoints.

---

## Version stamping (CI/CD)

The default version comes from [`VERSION`](./VERSION) and is exposed in code via `internal/version.Version`. Override at build time:

```sh
VERSION=$(cat VERSION)
go build -ldflags "-X github.com/Namularbre/knowledgeKeeperApi/internal/version.Version=${VERSION}" ./cmd/api
```

`/version` returns the stamped value, so CI/CD can tag images and verify a deployment matches the released version.

---

## Project layout

```
cmd/api/                 Composition root + swagger general annotations
internal/
  config/                Env loading
  version/               Build-stamped version variable
  infra/
    db/                  MariaDB connector + schema applier
    http/                Base HTTP server (meta routes, swagger mount, route registrar)
  auth/                  Auth bounded context (hexagonal)
    domain/              Entities, errors, ports (zero infra deps)
    app/                 Use cases (register, login, refresh)
    infra/               Adapters: bcrypt, JWT, MySQL repos, HTTP handlers
      sql/               Embedded bootstrap schema
docs/                    Generated OpenAPI artifacts (do not edit by hand)
scripts/swag.sh          Regenerate OpenAPI spec
```

---

## Common tasks

| Task | Command |
|---|---|
| Run full stack | `docker compose up --build` |
| Stop stack | `docker compose down` |
| Wipe DB | `docker compose down -v` |
| Run API only (against compose DB) | `docker compose up -d db && DB_HOST=127.0.0.1 go run ./cmd/api` |
| Regenerate OpenAPI | `./scripts/swag.sh` |
| Build production binary | `go build -ldflags "-s -w -X github.com/Namularbre/knowledgeKeeperApi/internal/version.Version=$(cat VERSION)" -o bin/api ./cmd/api` |
| Tidy dependencies | `go mod tidy` |
