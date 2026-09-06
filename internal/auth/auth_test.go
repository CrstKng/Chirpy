package auth

import (
	"testing"
	"time"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tokenSecret := "secret"
	ExpiresIn := time.Hour
	tests := []struct {
		name   string
		userID uuid.UUID
	}{
		{name: "case 1", userID: uuid.New()},
		{name: "case 2", userID: uuid.New()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := MakeJWT(tt.userID, tokenSecret, ExpiresIn)
			if err != nil {
				t.Fatalf("unexpected error creating token: %v", err)
			}
			if tokenString == "" {
				t.Error("expected non-empty token string")
			}
			gotID, err := ValidateJWT(tokenString, tokenSecret)
			if err != nil {
				t.Fatalf("unexpected error validating token: %v", err)
			}
			if gotID != tt.userID {
				t.Errorf("got %v, want %v", gotID, tt.userID)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	// setup: maybe create a valid token once here, or per test case
	userId := uuid.New()
	tests := []struct {
	name          string
	signingSecret string
	validateSecret string
	expiresIn     time.Duration
	wantErr       bool
}{
	{
		name:          "valid token",
		signingSecret: "secret-a",
		validateSecret: "secret-a",
		expiresIn:     time.Hour,
		wantErr:       false,
	},
	{
		name:          "wrong secret",
		signingSecret: "secret-a",
		validateSecret: "secret-b",
		expiresIn:     time.Hour,
		wantErr:       true,
	},
	{
		name:          "expired token",
		signingSecret: "secret-a",
		validateSecret: "secret-a",
		expiresIn:     -time.Hour,
		wantErr:       true,
	},
}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := MakeJWT(userId, tt.signingSecret, tt.expiresIn)
			if err != nil {
				t.Fatalf("unexpected error creating token: %v", err)
			}

			gotID, err := ValidateJWT(tokenString, tt.validateSecret)
			// now compare err against tt.wantErr, and gotID against userId if no error expected
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && gotID != userId {
				t.Errorf("ValidateJWT() gotID = %v, want %v", gotID, userId)
			}
		})
	}
}