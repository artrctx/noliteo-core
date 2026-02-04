package auth

import "github.com/lestrrat-go/jwx/v3/jwk"

type jwkManager struct {
	set jwk.Set
}

func NewJwkManager() jwkManager {
	return jwkManager{}
}
