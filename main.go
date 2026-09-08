package main

import (
	_ "github.com/lib/pq"
	"net/http"
	"time"
	"log"
	"sync/atomic"
	"os"
	"github.com/joho/godotenv"
	"database/sql"
	"github.com/CrstKng/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	PLATFORM string
	SECRET string
	POLKA_KEY string
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

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("error when connectiong to database: %s", err)
	}
	dbQueries := database.New(db)

	var apiConfig apiConfig
	apiCfg := &apiConfig
	apiCfg.dbQueries = dbQueries
	apiCfg.PLATFORM = os.Getenv("PLATFORM")
	apiCfg.SECRET = os.Getenv("SECRET")
	apiCfg.POLKA_KEY = os.Getenv("POLKA_KEY")

	mux := http.NewServeMux()

	filePath := http.Dir(filePathRoot)
	handler := http.FileServer(filePath)
	strippedFileServer := http.StripPrefix("/app", handler)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(strippedFileServer))

	mux.HandleFunc("GET /api/healthz", handlerReadinessCheck)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerRequestCounter)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUserEmailPassword)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerWebhooks)
	

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}