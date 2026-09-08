package main

import (
	"time"
	"encoding/json"
	"log"
	"net/http"

	"github.com/CrstKng/Chirpy/internal/auth"
	"github.com/CrstKng/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
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
	accessExpiresIn := time.Hour
	accessToken, err := auth.MakeJWT(user.ID, cfg.SECRET, accessExpiresIn)
	if err != nil {
		log.Printf("error when creating access token: %s", err)
		w.WriteHeader(500)
		return
	}
	refreshToken := auth.MakeRefreshToken()
	arg := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	}
	rToken, err := cfg.dbQueries.CreateRefreshToken(r.Context(), arg)
	if err != nil {
		log.Printf("error when creating refresh token: %s", err)
		w.WriteHeader(500)
		return
	}

	rVals = User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: accessToken,
		RefreshToken: rToken.Token,
		IsChirpyRed: user.IsChirpyRed,
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