package auth

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"log"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)

func HashPassword(password string) (string, error) {
	hashed_pass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return hashed_pass, fmt.Errorf("error when creating hash: %s\n", err)
	}
	return hashed_pass, nil
}


func CheckPasswordHash(password, hash string) (bool, error) {
	matches, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return matches, fmt.Errorf("error when comparing unhashed and hashed password: %s\n", err)
	}
	return matches, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:  "chirpy-access",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return signedToken, fmt.Errorf("error when signing the token: %s\n", err)
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var uuidNuLL uuid.UUID
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims , func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("error when parsing with claims token: %s", err)
		return uuidNuLL, err
	}
	stringID, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("error when getting subject from token.Claims: %s", err)
		return uuidNuLL, err
	}
	userID, err := uuid.Parse(stringID)
	if err != nil {
		log.Printf("error when converting userID from uuid.UUID to string: %s", err)
		return uuidNuLL, err
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	headerContent:= headers.Get("Authorization")
	if headerContent == "" {
		return "", fmt.Errorf("no authorization header")
	}
	splitAuthHeader := strings.Split(headerContent, " ")
	if len(splitAuthHeader) == 1 {
		return "", fmt.Errorf("invalid authorization header")
	}
	if splitAuthHeader[0] == "Bearer" {
		return splitAuthHeader[1], nil
	} else {
		return "", fmt.Errorf("authorization header is not written in the form: 'Bearer AuthString'")
	}
}