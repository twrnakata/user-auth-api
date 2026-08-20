# user-auth-api

Golang User Management API for the [7Solutions backend challenge](https://github.com/7-solutions/backend-challenge): MongoDB persistence, JWT HS256 auth, and a background job that logs the user count.

Lottery search is **design only** (no code):

- English: [`docs/lottery-search.md`](docs/lottery-search.md)
- Thai: [`docs/lottery-search.th.md`](docs/lottery-search.th.md)

## Stack

- Go 1.25
- Fiber
- MongoDB 7 (official driver v2)
- JWT HS256 (`github.com/golang-jwt/jwt/v5`)
- bcrypt for passwords

## Setup

1. Copy env:

```bash
cp .env.example .env
```

`.env` is used by Compose. Keep `MONGO_USERNAME` / `MONGO_PASSWORD` in sync with the URI if you run the API on the host.

### Option A — everything in Docker

```bash
docker compose up --build
```

Starts both Mongo and the API. No service name needed.

- API: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger`
- Mongo: `localhost:27017`

Inside the API container, Mongo is reached at hostname `mongo` (Compose overrides `MONGO_URI`). Wait until Mongo is healthy before the API starts.

Stop with `Ctrl+C`, or `docker compose down`.

### Option B — API on the host, Mongo in Docker only

```bash
docker compose up -d mongo
go run ./cmd/api
```

Name `mongo` on purpose so Compose does not also start the API container (that would fight `go run` for port `8080`).

Uses `MONGO_URI` from `.env` (`127.0.0.1`). Listens on `PORT` (default `8080`).

Swagger: `http://localhost:8080/swagger`  
Spec: `/swagger/openapi.yaml`

Specs are embedded in the API binary, so the same URLs work for both options.

| Variable | Purpose |
|---|---|
| `PORT` | HTTP listen port |
| `MONGO_URI` | Mongo connection string |
| `MONGO_DATABASE` | Database name (`user_management`) |
| `MONGO_USERNAME` / `MONGO_PASSWORD` | Used by `docker-compose.yml` to create the root user |
| `JWT_SECRET` | HMAC secret for HS256 (required) |

## Tests

```bash
go test ./...
```

Mongo repositories are unit-tested with fakes (no live database).

## JWT guide

`/auth/register` and `/auth/login` are public. Every `/users` route requires:

```http
Authorization: Bearer <token>
```

### How to get a token

1. `POST /auth/register` with `name`, `email`, `password`
2. `POST /auth/login` with the same `email` and `password`
3. Use `data.token` from the login response

The token is HS256, expires in **7 days**, and carries:

- `userId` — Mongo ObjectID hex
- `name`
- standard registered claims (`sub`, `exp`, `iat`, `nbf`)

Any valid token can call `/users` (there is no per-resource owner check).

Missing / malformed / expired tokens return **401**.

## Response envelope

All API routes except `GET /health` use:

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "errors": "optional",
  "serverTime": "2026-08-19T01:05:00+07:00"
}
```

| `code` | HTTP | Meaning |
|---|---|---|
| `0` | 200 / 201 | success |
| `400` | 400 | invalid parameter |
| `401` | 401 | unauthorized |
| `404` | 404 | not found |
| `409` | 409 | conflict (duplicate email) |
| `500` | 500 | internal error |

`createdAt` in `data` is `YYYY-MM-DD HH:mm:ss` in **Asia/Bangkok**. `serverTime` is RFC3339.

`GET /health` returns `{ "status": "ok" }` (no envelope).

## API samples

Base URL: `http://localhost:8080`

### Register

`POST /auth/register`

```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "password": "secret"
}
```

**201**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "69b3c0a1f2e3d4c5b6a79801",
    "name": "Alice",
    "email": "alice@example.com",
    "createdAt": "2026-08-19 01:05:00"
  },
  "serverTime": "2026-08-19T01:05:00+07:00"
}
```

Duplicate email → **409**. Missing fields / invalid JSON → **400**.

### Login

`POST /auth/login`

```json
{
  "email": "alice@example.com",
  "password": "secret"
}
```

**200**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "69b3c0a1f2e3d4c5b6a79801",
      "name": "Alice"
    }
  },
  "serverTime": "2026-08-19T01:05:00+07:00"
}
```

Unknown email and wrong password both return **401** `invalid credentials` (no user enumeration).

### List users

`GET /users`  
Header: `Authorization: Bearer <token>`

**200** — `data.users` is always an array (empty when there are no users).

### Get user

`GET /users/:id`

**200** — same user object as register (`id`, `name`, `email`, `createdAt`).  
Invalid id format → **400** `invalid user id`. Missing user → **404** `user not found`.

### Update user

`PUT /users/:id`

Send at least one of `name` or `email` (partial update):

```json
{
  "name": "Alice Updated"
}
```

**200** returns the user after update. Duplicate email → **409**. Empty body / invalid id format → **400**. Missing user → **404**.

### Delete user

`DELETE /users/:id`

**200** returns the deleted user. Invalid id format → **400**. Missing user → **404**.

## Logging and background job

Each request (except `/health`) writes one JSON access log: method, path, status, duration, `userId`, `requestId`.

Send `X-Request-ID` to correlate; if omitted the API generates a UUID and echoes it on the response.

Every **10 seconds** a goroutine counts documents in `users` and logs:

```json
{"timestamp":"2026-08-19 01:05:00","event":"userCount","count":3}
```

Unexpected 500s log the real error with `event: internalError` and `requestId`. The HTTP body stays generic (`internal server error`).

## Design decisions

- **Register is create.** There is no `POST /users`. The challenge asks for create + registration; one public register endpoint covers both.
- **Passwords never leave the API.** Stored as bcrypt hashes. User JSON has no password field.
- **Email unique index** on collection `users` in database `user_management`. Comparison is case-sensitive.
- **Invalid ObjectID is 400** (`invalid user id`). A well-formed id that is not in the database is **404** (`user not found`). Handler checks the 24-hex format before calling the service; the repository still converts with `ObjectIDFromHex` as a safety net.
- **PUT is a patch.** Empty string after trim is omitted; only non-empty `name` / `email` are `$set`.
- **No owner ACL.** JWT proves identity; any authenticated caller can list/update/delete any user. The challenge does not require ownership.
- **Hexagonal-ish layout:** domain ports, service, Mongo adapters, HTTP handlers. Handlers validate; services trust that input (except nil deps / nil output pointers).
- **Email format** is checked on register, login, and update (when email is sent) with the same pattern as iCRM (`test@google.com`).
- **Graceful shutdown:** `SIGINT` / `SIGTERM` stop HTTP first (`ShutdownWithTimeout` 10s is a maximum wait for in-flight requests, not a sleep; idle servers return immediately). Then the user-count job, then Mongo. `ReadTimeout` 30s only limits idle keepalive connections. gRPC is not implemented.

## Layout

```
cmd/api/                 entrypoint
docs/swagger/            Swagger UI + OpenAPI, embedded
docs/postman/            collection and local environment
docs/lottery-search.md     lottery design (English)
docs/lottery-search.th.md  lottery design (Thai)
internal/domain/         ports and domain errors
internal/service/        use cases
internal/repository/     Mongo adapters
internal/http/           Fiber handlers and routes
internal/middleware/     logging, recover, JWT
internal/job/            user-count ticker
pkg/                     JWT, Mongo, config, envelope, shared errors
```

## Postman

Import [`docs/postman/Backend_Challenge.postman_collection.json`](docs/postman/Backend_Challenge.postman_collection.json) and [`docs/postman/Backend_Challenge_Local.postman_environment.json`](docs/postman/Backend_Challenge_Local.postman_environment.json). Happy path: Health → Register → Login → `/users` with the saved token. More detail in [`docs/postman/README.md`](docs/postman/README.md).

## Swagger

Open `http://localhost:8080/swagger` after the API starts (local or Compose). Authorize with the JWT from `/auth/login`, then call `/users`.
