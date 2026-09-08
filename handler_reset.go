package main

import (
	"sync/atomic"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.PLATFORM != "dev" {
		log.Printf("Forbidden")
		w.WriteHeader(403)
		return
	}
	cfg.dbQueries.DeleteUsers(r.Context())
	var zeroValue atomic.Int32
	cfg.fileserverHits = zeroValue
	w.WriteHeader(200)
}