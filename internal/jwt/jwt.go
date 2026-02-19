package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/artrctx/noliteo-core/internal/config"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type JwtManager struct {
	cfg config.JwtConfig
}

var jwtMgr *JwtManager

// https://medium.com/techverito/secure-jwt-authentication-in-go-using-jwks-cba89d442f77
type Token struct {
	TID   uuid.UUID `json:"tid"`
	Ident string    `json:"ident"`
}

func newWithEnv() *JwtManager {
	return &JwtManager{config.GetJwtConfigFromEnv()}
}

func GenerateToken(t Token) (string, error) {
	if jwtMgr == nil {
		jwtMgr = newWithEnv()
	}

	token := jwt.New()

	tokenKeys := map[string]interface{}{
		"tid":           t.TID,
		"ident":         t.Ident,
		jwt.IssuedAtKey: time.Now(),
	}

	for k, v := range tokenKeys {
		if err := token.Set(k, v); err != nil {
			return "", fmt.Errorf("failed to set token key: %w", err)
		}
	}

	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), jwtMgr.cfg.Private))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return string(signedToken), nil
}

func ValidateToken(tkn string) (Token, error) {
	if jwtMgr == nil {
		jwtMgr = newWithEnv()
	}

	verifiedToken, err := jwt.Parse([]byte(tkn), jwt.WithKey(jwa.RS256(), jwtMgr.cfg.Public))
	if err != nil {
		return Token{}, fmt.Errorf("validating jwt failed: %w", err)
	}

	var tidStr, ident string
	verifiedToken.Get("tid", &tidStr)
	verifiedToken.Get("ident", &ident)

	if tidStr == "" {
		return Token{}, fmt.Errorf("jwt claim contains no tid (token id)")
	}
	tid, err := uuid.Parse(tidStr)
	if err != nil {
		return Token{}, fmt.Errorf("invalid format of jwt claim tid: %w", err)
	}

	if ident == "" {
		return Token{}, fmt.Errorf("jwt claim contains no identity")
	}

	return Token{tid, ident}, nil
}

func ValidateTokenFromRequest(r *http.Request) (Token, error) {
	var key string
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		splitKey := strings.Split(authHeader, " ")
		if len(splitKey) != 2 {
			return Token{}, fmt.Errorf("authorization header contains invalid key structure")
		}
		key = splitKey[1]
	} else {
		providedToken := r.URL.Query().Get("token")
		if providedToken == "" {
			return Token{}, fmt.Errorf("no token provided")
		}
		key = providedToken
	}

	token, err := ValidateToken(key)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}
