# TaskFlow — a hands-on Go course

A guided, project-based course that takes you from "I know Go basics" to
**shipping a production-ready Go service**. You build a minimal Jira-like task
manager, one phase at a time. You write the code; your mentor reviews each phase
before you move on.

The HTTP layer is built on the **[Gin](https://gin-gonic.com/) web framework**.

> Module path: `github.com/wchabir/taskflow` — rename later with
> `go mod edit -module <new>` plus a find/replace of imports if you push to your
> own repo.

## How this course works

Each phase follows the same loop:

1. **Goal & benefit** — what you'll build and why it matters in real projects.
2. **Hints & resources** — direct pointers (docs, packages, patterns), not full
   solutions. You do the thinking.
   - **Phases 1–3** (you're still new): hints are given as concrete
     *"go to…, write…"* steps, **followed by a full explanation of why each
     change was made**.
   - **Phases 4–10**: lighter-touch hints + resources + a definition of done.
3. **You build it** — write the code and create files yourself.
4. **Code review** — your mentor reviews the diff, tests edge cases, and
   suggests improvements. Only then do you advance.

Every phase's details are recorded in the **[Phase log](#phase-log)** at the
bottom of this file as you reach it.

## Prerequisites

- Go 1.26+ (`go version`) — installed ✓
- A terminal and an editor
- Docker — needed from **Phase 9** (install before then). *Not yet installed on
  this machine.*
- `golang-migrate` CLI — needed from **Phase 4** (we'll install it then).

## The project

A REST API for managing tasks (think a stripped-down Jira):

- Built on the **Gin** web framework
- Create / read / update / delete tasks
- Each task has a **status** and **priority** governed by `workflow.yaml`
- Protected by **Keycloak** (OIDC) authentication
- Backed by **PostgreSQL**
- A tiny vanilla-JS frontend at `/` (`static/index.html`) to test it live

## Syllabus

| Phase | Topic | You'll learn |
|------:|-------|--------------|
| **0** | Project scaffolding *(done for you)* | Layout, modules, tooling, compose, workflow file |
| **1** | HTTP server foundations (**Gin**) | Gin engine & routing, `/health`, JSON response envelope, graceful shutdown |
| **2** | Config & structured logging | Viper (`.env` + env vars), typed config, `slog`, loading `workflow.yaml` |
| **3** | Domain, validation & in-memory repo | Domain types, repository interface, dependency inversion, Gin binding + `go-playground/validator`, CRUD |
| **4** | PostgreSQL | `pgxpool`, `golang-migrate`, a real repository implementation |
| **5** | Service layer & business rules | Separation of concerns, workflow/transition enforcement, domain errors → HTTP |
| **6** | Testing | Table-driven tests, `testify`, mocks, `httptest`, integration tests |
| **7** | Authentication (Keycloak) | OIDC/JWT **Gin** middleware, `user_claims` in context, PKCE login in the frontend |
| **8** | Concurrency & goroutines | Background workers, channels, `context` cancellation, `errgroup`, worker pools |
| **9** | Docker | Multi-stage Dockerfile, full `docker-compose` stack, running migrations in containers |
| **10** | Production hardening | **Gin** middleware chain, request IDs, recovery, CORS, rate limiting, pagination, observability |

## Project layout

```
cmd/api/main.go        entry point — wires deps, starts the server
internal/
  config/              Viper config (.env + env) and workflow.yaml loader
  database/            pgxpool connection + ping
  domain/              core entities (Task, User) and value types
  handler/             Gin HTTP handlers
  middleware/          Gin middleware: request ID, logging, recovery, Keycloak OIDC
  repository/          storage interfaces + implementations
  service/             business logic
    mocks/             generated testify mocks
  validator/           go-playground/validator wrapper
migrations/            NNNNNN_name.{up,down}.sql (golang-migrate)
pkg/
  logger/              slog JSON (prod) / text (dev)
  response/            OK / Error / ValidationError envelope (Gin-aware)
static/index.html      vanilla-JS PKCE frontend (served at GET /)
docker/                docker-compose.yml (postgres + keycloak)
workflow.yaml          valid task statuses, priorities, transitions
.env                   local secrets (git-ignored; template in .env.example)
Makefile               developer task runner — `make help`
```

## Quick start

```bash
make help     # list all developer tasks
make run      # run the API
make build    # compile the binary into ./bin
make fmt vet  # format and vet
```

Open `http://localhost:8080/` against the running server once Phase 1 adds
`/health` — the status dot turns green.

---

# Phase log

A running record of each phase you complete, with its goal, the steps, and the
reasoning. Newest phase last.

## Phase 0 — Project scaffolding ✅

Done for you. Created the module (`github.com/wchabir/taskflow`, Go 1.26), the
full directory tree (each Go package has a `doc.go` describing its job), the
`Makefile`, `.env`/`.env.example`, `.gitignore`, `docker/docker-compose.yml`
(postgres + keycloak), `workflow.yaml` (statuses/priorities/transitions), and
the `static/index.html` frontend shell that pings `GET /health`. Verified with
`go build`, `go vet`, `gofmt`, and `make help`.

## Phase 1 — HTTP server foundations (Gin) 🚧

**Goal.** Replace the placeholder in `cmd/api/main.go` with a real **Gin** HTTP
server that (1) serves `GET /health` → `{"data":{"status":"ok"}}`, (2) serves
the frontend at `GET /`, (3) shuts down gracefully on Ctrl+C, and (4) writes all
responses through a reusable envelope in `pkg/response`.

**Benefit.** Every Go web service is a router + handlers + a shutdown story. Gin
is the most widely used Go web framework, so you'll learn the tool real teams
use — while still meeting graceful shutdown, the thing that stops a deploy from
dropping in-flight requests.

### Steps

1. **Go to the project root and add Gin.** Run `go get github.com/gin-gonic/gin`
   then `go mod tidy`. This downloads Gin and records it in `go.mod`/`go.sum`.
2. **Go to `pkg/response/response.go`** (create it next to `doc.go`). Write:
   - an envelope type: `type Envelope struct { Data any `json:"data,omitempty"`; Error *APIError `json:"error,omitempty"` }`
     and `type APIError struct { Message string `json:"message"` }`.
   - `func OK(c *gin.Context, status int, data any)` → calls `c.JSON(status, Envelope{Data: data})`.
   - `func Error(c *gin.Context, status int, msg string)` → calls `c.JSON(status, Envelope{Error: &APIError{Message: msg}})`.
3. **Go to `internal/handler/health.go`** (create it). Write
   `func Health(c *gin.Context)` that calls
   `response.OK(c, http.StatusOK, gin.H{"status": "ok"})`.
4. **Go to `cmd/api/main.go`** and rebuild `main()`:
   - create the router: `r := gin.Default()`.
   - register routes: `r.GET("/health", handler.Health)` and serve the frontend
     with `r.StaticFile("/", "./static/index.html")`.
   - build a server you control: `srv := &http.Server{Addr: ":8080", Handler: r}`.
5. **Add graceful shutdown** in `main()`:
   - start the server in a goroutine: `go srv.ListenAndServe()` (log the error
     unless it's `http.ErrServerClosed`).
   - block on signals: `ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)`; `defer stop()`; `<-ctx.Done()`.
   - drain: build a timeout context and call `srv.Shutdown(shutdownCtx)`.
6. **Run and test.** `make run`, then `curl -i localhost:8080/health`, open
   `http://localhost:8080/` (dot turns green), and press Ctrl+C to see a clean
   shutdown.

### Why each change

- **Why `go get gin`** — Gin is a thin, fast layer over `net/http`: a radix-tree
  router (path params, route groups), built-in JSON binding/validation, and a
  middleware chain. You asked to build on it, and it's the de-facto standard.
- **Why a `response` envelope package** — every endpoint returning the *same*
  JSON shape (`{"data":...}` / `{"error":...}`) gives the frontend one contract
  to code against and one place to change response format later. Wrapping
  `c.JSON` in `OK`/`Error` keeps handlers short and consistent, and prevents
  each handler from inventing its own shape. The frontend already reads
  `body.data.status`, which is why `/health` must nest under `data`.
- **Why a separate `handler` package** — keeping handlers out of `main` is the
  start of clean architecture: handlers become independently testable and `main`
  stays a thin wiring layer. You'll add many handlers in later phases.
- **Why `gin.Default()`** — it pre-installs two middlewares: a request *Logger*
  and a panic *Recovery* (so one bad handler doesn't crash the process). That's
  exactly what you want while learning; in Phase 10 you'll swap to `gin.New()`
  and attach your own middleware chain.
- **Why wrap Gin in `&http.Server{}` instead of `r.Run()`** — `r.Run()` is a
  convenience that calls `http.ListenAndServe` and hands you *no* server handle.
  To shut down gracefully you need the `*http.Server` value so you can call its
  `Shutdown` method. This is the single most important architectural choice in
  the phase.
- **Why run `ListenAndServe` in a goroutine** — it blocks forever serving
  requests. If you called it directly in `main`, you'd never reach the
  signal-waiting code. The goroutine serves while `main` waits for a shutdown
  signal.
- **Why `signal.NotifyContext` + `Shutdown`** — when you deploy, the
  orchestrator (or Ctrl+C) sends `SIGINT`/`SIGTERM`. `NotifyContext` turns those
  signals into a cancelled `context`. `srv.Shutdown(ctx)` then stops accepting
  new connections but lets in-flight requests finish (up to your timeout) before
  returning — no dropped requests mid-deploy. The timeout context is the safety
  valve so a stuck request can't block shutdown forever.
- **Why `http.ErrServerClosed` is not an error** — `Shutdown` causes
  `ListenAndServe` to return exactly this sentinel. Treating it as fatal would
  log a scary error on every clean stop, so you check for and ignore it.

### Definition of done

1. `make run` starts the server on `:8080`.
2. `curl -i localhost:8080/health` → `200`, `Content-Type: application/json`,
   body `{"data":{"status":"ok"}}`.
3. `http://localhost:8080/` shows the frontend with a **green** dot.
4. `curl -i -X POST localhost:8080/health` → `404`/`405` (Gin returns `404` for
   an unregistered method+path); `curl -i localhost:8080/nope` → `404`.
5. Ctrl+C logs a shutdown message and exits cleanly — no panic, no dropped
   request.
6. `gofmt` clean, `go vet ./...` clean, and `Health` uses `response.OK`.

### Resources

- Gin quickstart — https://gin-gonic.com/docs/quickstart/
- Gin examples (graceful shutdown) — https://github.com/gin-gonic/examples/tree/master/graceful-shutdown
- `http.Server.Shutdown` — https://pkg.go.dev/net/http#Server.Shutdown
- `signal.NotifyContext` — https://pkg.go.dev/os/signal#NotifyContext
