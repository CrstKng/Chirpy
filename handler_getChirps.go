package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/google/uuid"
	"sort"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	rVals := []Chirp{}
	if authorID != "" {
		parsedAuthorID, err := uuid.Parse(authorID)
		if err != nil {
			log.Printf("failed to parse authorID: %s", err)
			w.WriteHeader(500)
			return
		}
		chirpIDs, err := cfg.dbQueries.GetChirpIDByUserID(r.Context(), parsedAuthorID)
		if err != nil {
			log.Printf("author_id invalid or author_id has no chirps")
			w.WriteHeader(404)
			return
		}
		for _, chirpID := range chirpIDs {
			chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
			if err != nil {
				log.Printf("error when getting chirp from database: %s", err)
				w.WriteHeader(500)
				return
			}
			current := Chirp{
				ID: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserID: chirp.UserID,
			}
			rVals = append(rVals, current)	
		}
	} else {
		chirps, err := cfg.dbQueries.GetChirps(r.Context())
		if err != nil {
			log.Printf("error when getting chirps from database: %s", err)
			w.WriteHeader(500)
			return
		}
		for _, chirp := range chirps {
			current := Chirp{
				ID: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserID: chirp.UserID,
			}
			rVals = append(rVals, current)
		}
	}
	sortMethod := r.URL.Query().Get("sort")
	if sortMethod == "desc" {
		sort.Slice(rVals, func(i, j int) bool { return rVals[i].CreatedAt.After(rVals[j].CreatedAt) })
	} else {
		sort.Slice(rVals, func(i, j int) bool { return rVals[i].CreatedAt.Before(rVals[j].CreatedAt) })
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