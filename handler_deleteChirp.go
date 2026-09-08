package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/google/uuid"

	"github.com/CrstKng/Chirpy/internal/database"
	"github.com/CrstKng/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Authentication token JWT is malformed or not valid: %s", err)
		w.WriteHeader(401)
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.SECRET)
	if err != nil {
		log.Printf("Authentication token JWT not valid: %s", err)
		w.WriteHeader(403)
		return
	}
	chirpUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Invalid id: %s", err)
		w.WriteHeader(400)
		return
	}
	arg := database.ValidateChirpOwnerParams{
		ID: userID,
		ID_2: chirpUUID,
	}
	user, err := cfg.dbQueries.ValidateChirpOwner(r.Context(), arg)
	if err != nil {
		log.Printf("user not authorized to delete other user's chirps: %s", err)
		w.WriteHeader(403)
		return
	}
	err = cfg.dbQueries.DeleteChirpByID(r.Context(), chirpUUID)
	if err != nil {
		log.Printf("chirp not found: %s", err)
		w.WriteHeader(404)
		return
	}
	rVals := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	data, err := json.Marshal(rVals)
	if err != nil {
		log.Printf("error when marshaling JSON data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
	w.Write(data)
}