package main

import (
	"encoding/json"
	"log"
	"strings"
	"github.com/google/uuid"
	"time"
	"net/http"

	"github.com/CrstKng/Chirpy/internal/database"
	"github.com/CrstKng/Chirpy/internal/auth"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	type Parameters struct {
		Body   string `json:"body"`
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
	rVals := Chirp{}
	w.Header().Set("Content-Type", "application/json")

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error when getting bearer token: %s", err)
		w.WriteHeader(401)
		return
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.SECRET)
	if err != nil {
		log.Printf("Authentication token JWT not valid: %s", err)
		w.WriteHeader(401)
		return
	}

	dbParams := database.CreateChirpParams{
		Body: params.Body,
		UserID: userID,
	}

	ctx := r.Context()
	chirp, err := cfg.dbQueries.CreateChirp(ctx, dbParams)
	if err != nil {
		log.Printf("error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	rVals = Chirp{
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
