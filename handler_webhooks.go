package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/google/uuid"

	"github.com/CrstKng/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.POLKA_KEY {
		log.Printf("incorrect api key")
		w.WriteHeader(401)
		return
	}
	type data_type struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data data_type `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if params.Event != "user.upgraded" {
		log.Printf("incorrect event: %s", err)
		w.WriteHeader(204)
		return
	}
	err = cfg.dbQueries.UpdateUserMembership(r.Context(), params.Data.UserID)
	if err != nil {
		log.Printf("user is not found: %s", err)
		w.WriteHeader(404)
		return
	}
	updated_user, err := cfg.dbQueries.GetUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		log.Printf("error when getting user data by id: %s", err)
		w.WriteHeader(500)
		return
	}
	rVals := User{
		ID: updated_user.ID,
		CreatedAt: updated_user.CreatedAt,
		UpdatedAt: updated_user.UpdatedAt,
		Email: updated_user.Email,
		IsChirpyRed: updated_user.IsChirpyRed,
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