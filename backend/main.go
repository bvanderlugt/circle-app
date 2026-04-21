package main

import (
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

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))

	srv := server.New(mux, nil)
	log.Printf("serving %s on %s", staticDir, addr)
	if err := srv.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
