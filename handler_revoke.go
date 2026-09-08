package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CrstKng/Chirpy/internal/auth"
)

type RefreshToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Authentication token JWT is malformed or not valid: %s", err)
		w.WriteHeader(401)
		return
	}
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("error when revoking refresh token %s: %s", refreshToken, err)
		w.WriteHeader(401)
		return
	}
	rVals := RefreshToken{Token: refreshToken}
	data, err := json.Marshal(rVals)
	if err != nil {
		log.Printf("error when marshaling JSON data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
	w.Write(data)
}