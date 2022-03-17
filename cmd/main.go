package main

import (
	"log"
	"net/http"

	"github.com/zfirdavs/devim-lab/internal/coordinates"
)

func main() {
	mux := http.NewServeMux()

	coordsHandler := coordinates.NewHandler()
	mux.HandleFunc("/", coordsHandler.GetCoordinates())

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("server is listening on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
