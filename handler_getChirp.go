package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Invalid id: %s", err)
		w.WriteHeader(400)
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(r.Context(), uuid)
	if err != nil {
		log.Printf("Chirp with id: %v not found: %s", r.URL, err)
		w.WriteHeader(404)
		return
	}
	
	rVals := Chirp{
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
	_, err = w.Write(data)
	if err != nil {
		log.Printf("error when writing data to body: %s", err)
	}
	w.WriteHeader(200)
}