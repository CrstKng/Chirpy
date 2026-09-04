package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
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