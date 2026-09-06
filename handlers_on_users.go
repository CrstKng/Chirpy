package main

import (
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"time"
	"net/http"

	"github.com/CrstKng/Chirpy/internal/database"
	"github.com/CrstKng/Chirpy/internal/auth"
)


type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}
	userParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	}
	ctx := r.Context()
	user, err := cfg.dbQueries.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	rVals := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	data, err := json.Marshal(rVals)
	if err != nil {
		log.Printf("error when marshaling JSON data: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(data)
}