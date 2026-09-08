package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/CrstKng/Chirpy/internal/auth"
)

type AccessToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Authentication token JWT is malformed or not valid: %s", err)
		w.WriteHeader(401)
		return
	}
	user, err := cfg.dbQueries.GetUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("error when getting user by refresh from db: %s", err)
		w.WriteHeader(401)
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, cfg.SECRET, time.Hour)
	if err != nil {
		log.Printf("error when creating access token: %s", err)
		w.WriteHeader(500)
		return
	}
	rVals := AccessToken{Token: accessToken}
	data, err := json.Marshal(rVals)
	if err != nil {
		log.Printf("error when marshaling JSON data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}