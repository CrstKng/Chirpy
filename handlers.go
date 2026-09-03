package main

import (
	"net/http"
	"fmt"
)

func handlerReadinessCheck(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)
	res := "OK"
	_, err := resWriter.Write([]byte(res))
	if err != nil {
		fmt.Printf("error when writing data to body: %s\n", err)
	}
}