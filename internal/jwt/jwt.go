package jwt

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/artrctx/noliteo-core/internal/config"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type JwtManager struct {
	cfg config.JwtConfig
}

var jwtMgr *JwtManager

// https://medium.com/techverito/secure-jwt-authentication-in-go-using-jwks-cba89d442f77
type Token struct {
	TID   string `json:"tid"`
	Ident string `json:"ident"`
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

func VerifyToken(tkn string) (Token, error) {
	if jwtMgr == nil {
		jwtMgr = newWithEnv()
	}

	verifiedToken, err := jwt.Parse([]byte(tkn), jwt.WithKey(jwa.RS256(), jwtMgr.cfg.Public))
	if err != nil {
		return Token{}, fmt.Errorf("verifying jwt failed: %w", err)
	}

	var tid, ident string
	verifiedToken.Get("tid", &tid)
	verifiedToken.Get("ident", &ident)

	if tid == "" {
		return Token{}, fmt.Errorf("jwt claim contains no tid (token id)")
	}

	if ident == "" {
		return Token{}, fmt.Errorf("jwt claim contains no identity")
	}

	return Token{tid, ident}, nil
}

func VerifyTokenFromRequest(r *http.Request) (Token, error) {
	splitKey := strings.Split(r.Header.Get("Authorization"), " ")
	if len(splitKey) != 2 {
		return Token{}, fmt.Errorf("authorization header contains invalid key structure")
	}
	key := splitKey[1]

	token, err := VerifyToken(key)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}
