package main

import (
	"fmt"
	"sync/atomic"
	"log"
	"net/http"
)

func handlerReadinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	res := "OK"
	_, err := w.Write([]byte(res))
	if err != nil {
		log.Printf("error when writing data to body: %s", err)
	}
}

func (cfg *apiConfig) handlerRequestCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`,
	cfg.fileserverHits.Load())))
}

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