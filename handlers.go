package main

import (
	"net/http"
	"encoding/json"
	"fmt"
)

func ReadinessCheck(resWriter http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)
	res := "OK"
	data, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("error when marshaling response: %s\n", err)
	}
	_, err = resWriter.Write(data)
	if err != nil {
		fmt.Printf("error when writing data to body: %s\n", err)
	}
}