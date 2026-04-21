# circle-app

Go backend (`gocloud.dev/server`) serving a React (Vite) frontend from disk.

## Layout

```
circle-app/
  frontend/     # Vite + React app; `npm run build` → frontend/dist/
  backend/      # Go server; serves files from $STATIC_DIR (default ./static)
  Dockerfile    # multi-stage build → distroless image
```

## Local development

Build the frontend, then run the backend pointed at the built output:

```sh
# 1. build the frontend
cd frontend
npm install
npm run build

# 2. run the backend, pointing it at the built assets
cd ../backend
STATIC_DIR=../frontend/dist go run .
```

Open http://localhost:8080.

For an iterative frontend loop, run `npm run dev` in `frontend/` (port 5173)
while the Go server runs separately — the two are fully decoupled.

## Docker

```sh
docker build -t circle-app .
docker run --rm -p 8080:8080 circle-app
```

## Environment variables

| Var         | Default     | Purpose                          |
|-------------|-------------|----------------------------------|
| `ADDR`      | `:8080`     | Listen address                   |
| `STATIC_DIR`| `./static`  | Directory of built frontend files|
