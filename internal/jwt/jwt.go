package jwt

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Token struct {
	Ident string `json:"ident"`
}

func GenerateToken() (string, error) {
	token := jwt.New()

	tokenKeys := map[string]interface{}{
		jwt.IssuedAtKey: time.Now(),
	}

	for k, v := range tokenKeys {
		if err := token.Set(k, v); err != nil {
			return "", fmt.Errorf("failed to set token key with error: %w", err)
		}
	}

	return "", nil
}

func VerifyToken(token string) (Token, error) {
	return Token{}, nil
}
