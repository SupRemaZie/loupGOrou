package main

import (
	"flag"
	"log"

	"github.com/SupRemaZie/loupGOrou/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "adresse d'écoute")
	static := flag.String("static", "./web/dist", "dossier du frontend compilé")
	flag.Parse()

	log.Printf("Serveur démarré sur http://localhost%s", *addr)
	srv := server.New(*addr, *static)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
