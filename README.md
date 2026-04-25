# circle-app

Go backend (`gocloud.dev/server`) + Postgres, serving a React (Vite) frontend
from disk. Trivial CRUD on named "circles".

## Layout

```
circle-app/
  frontend/            # Vite + React app; `npm run build` → frontend/dist/
  backend/             # Go server; serves static files + /api
  Dockerfile           # multi-stage → distroless image
  docker-compose.yml   # db + app for local/dev/prod-like runs
```

## API

| Method | Path                 | Body / response                                  |
|--------|----------------------|--------------------------------------------------|
| GET    | `/api/circles/count` | `{"count": N}`                                   |
| GET    | `/api/circles`       | `[{"id":1,"name":"alpha","created_at":"..."}]`   |
| POST   | `/api/circles`       | body `{"name":"alpha"}` → circle JSON (201)      |
| DELETE | `/api/circles/{id}`  | 204, or 404 if missing                           |

Validation errors on POST return 400 with `{"error":"..."}`. Name rules: 1-50
characters after trimming, no control characters.

## Local development

### Option 1: full stack in Docker

```sh
docker compose up --build
# → http://localhost:8080
```

### Option 2: Postgres in Docker, backend + frontend locally

```sh
docker compose up -d db

# build the frontend once
cd frontend && npm install && npm run build

# run the backend
cd ../backend
DATABASE_URL=postgres://circle:circle@localhost:5432/circle?sslmode=disable \
STATIC_DIR=../frontend/dist go run .
```

For iterative frontend work, run `npm run dev` in `frontend/` (port 5173)
with a Vite proxy forwarding `/api` to `localhost:8080`.

### Option 3: no DB, in-memory store

Omit `DATABASE_URL` and the backend falls back to an in-memory store — handy
for UI iteration without Postgres running.

## Tests

```sh
cd backend && go test ./...
```

No database is required; tests cover validation logic only.

## Environment variables

| Var            | Default                              | Purpose                           |
|----------------|--------------------------------------|-----------------------------------|
| `ADDR`         | `:8080`                              | Listen address                    |
| `STATIC_DIR`   | `./static`                           | Built frontend directory          |
| `DATABASE_URL` | *(unset → in-memory store)*          | Postgres connection URL           |
