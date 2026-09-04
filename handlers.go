package main

import (
	"net/http"
	"fmt"
	"sync/atomic"
	"encoding/json"
	"log"
	"strings"
)

func handlerReadinessCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	res := "OK"
	_, err := w.Write([]byte(res))
	if err != nil {
		log.Printf("error when writing data to body: %s", err)
	}
}

func (cfg *apiConfig) handlerRequestCounter(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`,
	cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	var zeroValue atomic.Int32
	cfg.fileserverHits = zeroValue
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	type returnVals struct {
		Cleaned string `json:"cleaned_body"`
		Error string `json:"error"`
	}

	rVals := returnVals{}
	w.Header().Set("Content-Type", "application/json")

	if len(params.Body) <= 140 {
		splitString := strings.Split(params.Body, " ")
		for i, word := range splitString {
			loweredWord := strings.ToLower(word)
			for _, profanity := range profaneWords{
				if loweredWord == profanity {
					splitString[i] = "****"
				}
			}
		}
		cleanedString := strings.Join(splitString, " ")

		rVals.Cleaned = cleanedString
		w.WriteHeader(200)

	} else {
		rVals.Error = "Chirp is too long"
		w.WriteHeader(400)
	}

	data, err := json.Marshal(rVals)
	if err != nil {
		log.Printf("error when marshaling JSON data: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Write(data)
}