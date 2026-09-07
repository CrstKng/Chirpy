package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func MakeRefreshToken() string {
	byteData := make([]byte, 32)
	_, err := rand.Read(byteData)
	if err != nil {
		log.Printf("error when generating random byte date: %s", err)
	}
	return hex.EncodeToString(byteData)
}