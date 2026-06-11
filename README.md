# TaskFlow — a hands-on Go course

A guided, project-based course that takes you from "I know Go basics" to
**shipping a production-ready Go service**. You build a minimal Jira-like task
manager, one phase at a time. You write the code; your mentor reviews each phase
before you move on.

> Module path: `github.com/wchabir/taskflow` — rename later with
> `go mod edit -module <new>` plus a find/replace of imports if you push to your
> own repo.

## How this course works

Each phase follows the same loop:

1. **Goal & benefit** — what you'll build and why it matters in real projects.
2. **Hints & resources** — direct pointers (docs, packages, patterns), not full
   solutions. You do the thinking.
3. **You build it** — write the code and create files yourself.
4. **Code review** — your mentor reviews the diff, tests edge cases, and
   suggests improvements. Only then do you advance.

## Prerequisites

- Go 1.26+ (`go version`) — installed ✓
- A terminal and an editor
- Docker — needed from **Phase 9** (install before then). *Not yet installed on
  this machine.*
- `golang-migrate` CLI — needed from **Phase 4** (we'll install it then).

## The project

A REST API for managing tasks (think a stripped-down Jira):

- Create / read / update / delete tasks
- Each task has a **status** and **priority** governed by `workflow.yaml`
- Protected by **Keycloak** (OIDC) authentication
- Backed by **PostgreSQL**
- A tiny vanilla-JS frontend at `/` (`static/index.html`) to test it live

## Syllabus

| Phase | Topic | You'll learn |
|------:|-------|--------------|
| **0** | Project scaffolding *(done for you)* | Layout, modules, tooling, compose, workflow file |
| **1** | HTTP server foundations | `net/http`, `ServeMux` routing, `/health`, JSON response envelope, graceful shutdown |
| **2** | Config & structured logging | Viper (`.env` + env vars), typed config, `slog`, loading `workflow.yaml` |
| **3** | Domain, validation & in-memory repo | Domain types, repository interface, dependency inversion, `go-playground/validator`, CRUD |
| **4** | PostgreSQL | `pgxpool`, `golang-migrate`, a real repository implementation |
| **5** | Service layer & business rules | Separation of concerns, workflow/transition enforcement, domain errors → HTTP |
| **6** | Testing | Table-driven tests, `testify`, mocks, `httptest`, integration tests |
| **7** | Authentication (Keycloak) | OIDC/JWT middleware, `user_claims` in context, PKCE login in the frontend |
| **8** | Concurrency & goroutines | Background workers, channels, `context` cancellation, `errgroup`, worker pools |
| **9** | Docker | Multi-stage Dockerfile, full `docker-compose` stack, running migrations in containers |
| **10** | Production hardening | Middleware chain, request IDs, recovery, CORS, rate limiting, pagination, observability |

## Project layout

```
cmd/api/main.go        entry point — wires deps, starts the server
internal/
  config/              Viper config (.env + env) and workflow.yaml loader
  database/            pgxpool connection + ping
  domain/              core entities (Task, User) and value types
  handler/             HTTP handlers
  middleware/          request ID, logging, recovery, Keycloak OIDC
  repository/          storage interfaces + implementations
  service/             business logic
    mocks/             generated testify mocks
  validator/           go-playground/validator wrapper
migrations/            NNNNNN_name.{up,down}.sql (golang-migrate)
pkg/
  logger/              slog JSON (prod) / text (dev)
  response/            OK / Error / ValidationError envelope
static/index.html      vanilla-JS PKCE frontend (served at GET /)
docker/                docker-compose.yml (postgres + keycloak)
workflow.yaml          valid task statuses, priorities, transitions
.env                   local secrets (git-ignored; template in .env.example)
Makefile               developer task runner — `make help`
```

## Quick start

```bash
make help     # list all developer tasks
make run      # run the API (Phase 0: prints a placeholder banner)
make build    # compile the binary into ./bin
make fmt vet  # format and vet
```

Open `static/index.html` against the running server once Phase 1 adds `/health`.
