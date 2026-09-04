package main

import (
	"net/http"
	"fmt"
	"sync/atomic"
	"encoding/json"
	"log"
	"strings"
	"github.com/google/uuid"
	"time"
	"github.com/CrstKng/Chirpy/internal/database"
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

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	type returnVals struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}
	rVals := returnVals{}
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	user, err := cfg.dbQueries.CreateUser(ctx, params.Email)
	if err != nil {
		log.Printf("error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	rVals = returnVals{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	data, err := json.Marshal(rVals)
	if err != nil {
		log.Printf("error when marshaling JSON data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(data)
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	type Parameters struct {
		Body   string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) <= 140 {
		splitString := strings.Split(params.Body, " ")
		for i, word := range splitString {
			loweredWord := strings.ToLower(word)
			for _, profanity := range profaneWords{
				if loweredWord == profanity {
					splitString[i] = "****"
				}
			}
		}
		cleanedString := strings.Join(splitString, " ")
		params.Body = cleanedString
	} else {
		log.Printf("Chirp is too long")
		w.WriteHeader(400)
		return
	}

	type returnVals struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	rVals := returnVals{}
	w.Header().Set("Content-Type", "application/json")

	dbParams := database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
	}

	ctx := r.Context()
	chirp, err := cfg.dbQueries.CreateChirp(ctx, dbParams)
	if err != nil {
		log.Printf("error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	rVals = returnVals{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	data, err := json.Marshal(rVals)
	if err != nil {
		log.Printf("error when marshaling JSON data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(data)
}