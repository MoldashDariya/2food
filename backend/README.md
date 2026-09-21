# 2food backend

Go REST API for 2food, plus the background workers and database migrations. See the [main README](../README.md) for the overall idea and the mobile app.

## Run it

```bash
cp .env.example .env
docker compose up -d postgres redis rustfs      # add the MINIO_* variables to .env first, see the main README
go run ./cmd/migrator -up -user postgres -password postgres -dbname 2food
go run ./cmd/2food
```

The API is served at `http://localhost:8080/api/v1`, the interactive docs at `http://localhost:8080/docs`. Run the tests with `go test ./...`.

## Migrator

`cmd/migrator` wraps golang-migrate. Connection flags: `-host`, `-port`, `-user`, `-password`, `-dbname`, `-sslmode`, `-migrations-path`. Actions (pick one):

| Flag | Action |
| --- | --- |
| `-up` | Apply all pending migrations |
| `-down` | Roll back one step |
| `-drop` | Roll back everything |
| `-force N` | Force the version to N (for a dirty state) |

## Code layout

Every feature (`auth`, `user`, `restaurant`, `product`, `order`, `review`, `catalog`, `upload`) is a package with a `service` (business logic, works with small interfaces) and `transport/http` (handlers and DTOs). Entities and errors are in `internal/domain`, integrations (Postgres, Stripe, Twilio, SMTP, FCM, S3, JWT, queue) are in `internal/infra`, cron and queue handlers are in `internal/worker`. Everything is wired together by hand in `internal/app/app.go`, and the routes are in `internal/app/router.go`.

## Background jobs

The worker and the cron scheduler run inside the API process (asynq on top of Redis):

- payment confirm, cancel, refund and audit jobs
- push notification delivery
- cancelling orders that were never paid (`WORKER_ORDER_CLEANER_CRON`, `WORKER_PENDING_ORDER_TTL`)
- role cleanup (`WORKER_ROLE_CLEANER_CRON`): users with the partner role who no longer work at any restaurant go back to customer

## API docs

Swagger files in `docs/` are generated from the code annotations, do not edit them by hand. Regenerate with `swag init -g cmd/2food/main.go` (see `docs/README.md`).
