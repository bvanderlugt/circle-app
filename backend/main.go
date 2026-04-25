package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"gocloud.dev/server"
)

func main() {
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./static"
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	store, cleanup := newStore(context.Background())
	defer cleanup()

	mux := http.NewServeMux()
	registerAPIRoutes(mux, store)
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))

	srv := server.New(mux, nil)
	log.Printf("serving %s on %s", staticDir, addr)
	if err := srv.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func newStore(ctx context.Context) (Store, func()) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Printf("DATABASE_URL unset; using in-memory store")
		return NewMemStore(), func() {}
	}
	pg, err := NewPgStore(ctx, url)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	log.Printf("connected to postgres")
	return pg, func() { _ = pg.Close() }
}
