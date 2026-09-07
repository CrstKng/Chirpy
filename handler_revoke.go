package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type RefreshToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	splitRefreshHeader := strings.Split(r.Header["Authorization"][0], " ")
	if splitRefreshHeader[0] != "Bearer" {
		log.Printf("authorization header is not written in the form: 'Authorization: Bearer <refresh-token>'")
		w.WriteHeader(500)
		return
	}
	refreshToken := splitRefreshHeader[1]
	err := cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
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