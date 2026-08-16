# User Management API

A RESTful API in Go for managing users with a name and date of birth. Age is never stored — it is calculated from `dob` on every read using Go's `time` package.

## Tech stack

| | |
|---|---|
| HTTP | [Fiber](https://github.com/gofiber/fiber) v2 |
| Database | PostgreSQL 16 |
| Data access | [sqlc](https://sqlc.dev) — generated, no ORM |
| Logging | [Uber Zap](https://github.com/uber-go/zap) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |

## Quick start

Requires Docker and Docker Compose.

```bash
git clone https://github.com/adityaonrepeat/user-management-api.git
cd user-management-api
cp .env.example .env
docker compose up --build
```

The API listens on `http://localhost:3000`. The `users` table is created automatically on first start.

```bash
curl http://localhost:3000/health
# {"status":"ok"}
```

## Running locally

Requires Go 1.26+ and a reachable PostgreSQL instance.

```bash
docker compose up -d postgres    # database only
go run ./cmd/server
```

Configuration comes from the environment, with defaults matching `.env.example`:

| Variable | Default |
|---|---|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/user_management?sslmode=disable` |
| `SERVER_PORT` | `3000` |

## API

`dob` is always `YYYY-MM-DD`. Every response carries an `X-Request-ID` header.

### Create a user

```bash
curl -X POST http://localhost:3000/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```
```json
{"id": 1, "name": "Alice", "dob": "1990-05-10"}
```
`201 Created`

### Get a user

```bash
curl http://localhost:3000/users/1
```
```json
{"id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35}
```
`200 OK`

### List users

```bash
curl http://localhost:3000/users
```
```json
[{"id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35}]
```
`200 OK`

Optional pagination. Without query parameters every user is returned.

```bash
curl "http://localhost:3000/users?limit=10&offset=0"
```

`limit` must be between 1 and 100; `offset` must be non-negative. The total number of users is returned in the `X-Total-Count` header, keeping the response body a plain array.

### Update a user

```bash
curl -X PUT http://localhost:3000/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'
```
```json
{"id": 1, "name": "Alice Updated", "dob": "1991-03-15"}
```
`200 OK`

### Delete a user

```bash
curl -X DELETE http://localhost:3000/users/1
```
`204 No Content`

### Errors

| Status | When |
|---|---|
| `400` | Malformed JSON, failed validation, non-integer id, invalid pagination |
| `404` | No user with that id, including on `PUT` and `DELETE` |
| `500` | Unexpected failure |

```json
{"error": {"code": "NOT_FOUND", "message": "user not found"}}
```

> Age is calculated at request time, so returned ages may differ from the example values above depending on the current date.

## Project structure

```
cmd/server/        entrypoint and dependency wiring
config/            environment variables into a typed config
db/
  migrations/      schema as numbered SQL files
  queries/         hand-written SQL with sqlc annotations
  sqlc/            generated data access code
internal/
  handler/         HTTP: binding, validation, status codes
  repository/      thin wrapper over the generated queries
  service/         business rules and age calculation
  routes/          route registration
  middleware/      request id, request duration logging
  models/          domain types and request/response DTOs
  logger/          Zap construction
```

Requests flow in one direction: `handler → service → repository → sqlc`. The handler never touches the database and the service never imports Fiber.

## Design decisions

**Age is derived, never stored.** A stored age silently becomes wrong at midnight on the user's birthday. Deriving it from `dob` on every read means it cannot drift.

**`dob` is a `DATE`, not a `TIMESTAMP`.** A date of birth is a calendar fact, not a moment in time. `DATE` carries no time component and no timezone, so the stored value cannot shift across zone boundaries and age remains deterministic. Ages are computed in UTC for the same reason.

**Age is calculated by comparing month then day, not `YearDay()`.** In a leap year every date after February has a `YearDay` one higher than in a common year. That offset cancels out the day before a birthday and yields an age one year too high — someone born 1990-07-15 would be reported as 34 on 2024-07-14 rather than 33. The unit tests cover this case explicitly.

**The listing uses a single reference time.** Ages for all users in one response are computed from the same instant, so a request spanning midnight cannot report inconsistent ages.

**Response shapes match the specification exactly.** `POST` and `PUT` return `id`, `name` and `dob`; the `GET` endpoints additionally return `age`. Rather than impose one uniform response type, there are two, because the contract defines two.

**`UPDATE` and `DELETE` distinguish "missing" from "succeeded".** Updating or deleting a row that does not exist is not an SQL error — it simply affects no rows. `UpdateUser` uses `RETURNING` so a missing row surfaces as `ErrNoRows`, and `DeleteUser` reports its affected-row count. Both become `404` instead of a misleading `200` or `204`.

**Validation is split by concern.** Whether `dob` is a well-formed date is a transport question, handled by a validator tag. Whether it is a *usable* birth date — specifically, not in the future — is a business rule, handled in the service. A future date of birth produces a negative age, so rejecting it protects the one value this API exists to compute.

**The repository is deliberately thin.** sqlc already generates a `Querier`, so the repository exists for one reason: it is the boundary where generated and driver types stop. `pgx.ErrNoRows` and `db.User` live below it; `ErrNotFound` and `models.User` live above it. Without it, driver types would reach the service layer.

**Migrations are plain SQL with no migration tool.** The files follow golang-migrate's naming convention so the project can adopt it later without renaming anything, but adding a migration framework for a single table would be over-engineering.

**Generated code is committed.** The project builds from a clean clone without sqlc installed.

**Errors are rendered in one place.** A single Fiber `ErrorHandler` produces the error envelope, so no handler formats errors itself. Unmapped errors are logged with full detail and returned as a generic `500`, so internals never leak to clients.

## Testing

```bash
go test ./... -v
```

The age calculation is covered by a table-driven test including leap-year boundaries, birthdays on the current date, and the day either side of a birthday.

## Development

```bash
sqlc generate     # regenerate the data access layer after changing SQL
go vet ./...
gofmt -l .
```

Regenerating requires [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html); it is not needed to build or run the project.

## Possible extensions

Deliberately out of scope for this implementation, but natural next steps:

- Cursor-based pagination, which is stable under concurrent inserts in a way that offset pagination is not
- A partial-update `PATCH` endpoint alongside the full-replacement `PUT`
- Filtering and sorting on the listing
- Rate limiting and authentication
- Structured request tracing via OpenTelemetry, building on the existing request id
