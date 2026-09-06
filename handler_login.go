package main

import (
	"strconv"
	"time"
	"encoding/json"
	"log"

	"net/http"

	"github.com/CrstKng/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = 3600
	}
	
	rVals := User{}
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	user, err1 := cfg.dbQueries.GetUserByEmail(ctx, params.Email)
	matches, err2 := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err1 != nil || err2 != nil || !matches {
		log.Printf("incorrect email or password")
		w.WriteHeader(401)
		return
	}
	expiresIn, err := time.ParseDuration(strconv.Itoa(params.ExpiresInSeconds) + "s")
	if err != nil {
		log.Printf("error when converting expiresInSeconds to time.Duration: %s", err)
		w.WriteHeader(500)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.SECRET, expiresIn)
	if err != nil {
		log.Printf("error when creating token: %s", err)
		w.WriteHeader(500)
		return
	}
	rVals = User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
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