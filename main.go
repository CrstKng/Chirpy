package main

import (
	"net/http"
	"time"
	"log"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	var apiConfig apiConfig
	apiCfg := &apiConfig

	mux := http.NewServeMux()

	filePath := http.Dir(filePathRoot)
	handler := http.FileServer(filePath)
	strippedFileServer := http.StripPrefix("/app", handler)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(strippedFileServer))

	mux.HandleFunc("GET /api/healthz", handlerReadinessCheck)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerRequestCounter)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}