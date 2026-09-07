package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/CrstKng/Chirpy/internal/auth"
)

type AccessToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	splitRefreshHeader := strings.Split(r.Header["Authorization"][0], " ")
	if splitRefreshHeader[0] != "Bearer" {
		log.Printf("authorization header is not written in the form: 'Authorization: Bearer <refresh-token>'")
		w.WriteHeader(500)
		return
	}
	refreshToken := splitRefreshHeader[1]
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