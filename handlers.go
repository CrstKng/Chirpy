package main

import (
	"net/http"
	"fmt"
	"sync/atomic"
)

func handlerReadinessCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	res := "OK"
	_, err := w.Write([]byte(res))
	if err != nil {
		fmt.Printf("error when writing data to body: %s\n", err)
	}
}

func (cfg *apiConfig) handlerRequestCounter(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %d\n", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	var zeroValue atomic.Int32
	cfg.fileserverHits = zeroValue
}