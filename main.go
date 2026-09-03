package main

import (
	"net/http"
	"time"
	"log"
)

func main() {
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	filePath := http.Dir(filePathRoot)
	handler := http.FileServer(filePath)
	strippedFileServer := http.StripPrefix("/app", handler)
	mux.Handle("/app/", strippedFileServer)

	mux.HandleFunc("/healthz", ReadinessCheck)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}