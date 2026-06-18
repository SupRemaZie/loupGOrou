package main

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/SupRemaZie/loupGOrou/internal/server"
)

func main() {
	// Serve the built React app from web/dist (relative to this file's directory)
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..", "web", "dist")

	addr := ":8080"
	log.Printf("Serveur démarré sur http://localhost%s", addr)

	srv := server.New(addr, root)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
