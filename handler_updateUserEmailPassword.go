package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CrstKng/Chirpy/internal/database"
	"github.com/CrstKng/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpdateUserEmailPassword(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Authentication token JWT is malformed or not valid: %s", err)
		w.WriteHeader(401)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.SECRET)
	if err != nil {
		log.Printf("Authentication token JWT not valid: %s", err)
		w.WriteHeader(401)
		return
	}
	arg := database.UpdateUserEmailPasswordParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
		ID: userID,
	}
	err = cfg.dbQueries.UpdateUserEmailPassword(r.Context(), arg)
	if err != nil {
		log.Printf("error updating user's email and password: %s", err)
		w.WriteHeader(500)
		return
	}
	new_user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("error getting user's info by email: %s", err)
		w.WriteHeader(500)
		return
	}
	rVals := User{
		ID: new_user.ID,
		CreatedAt: new_user.CreatedAt,
		UpdatedAt: new_user.UpdatedAt,
		Email: new_user.Email,
		IsChirpyRed: new_user.IsChirpyRed,
	}
	data, err := json.Marshal(rVals)
	if err != nil {
		log.Printf("error when marshaling JSON data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}