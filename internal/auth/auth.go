package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/httprc/v3/tracesink"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

var (
	JwksUrl  = ("JWKS_URL")
	JwkCache *jwk.Cache
)

// errors
var (
	ErrMissingUserID = errors.New("missing user id")
)

func getCachedKeyset(ctx context.Context) (*jwk.Cache, error) {
	c, err := jwk.NewCache(ctx, httprc.NewClient(
		httprc.WithTraceSink(tracesink.NewSlog(slog.New(slog.NewJSONHandler(os.Stderr, nil)))),
	))

	if err != nil {
		return nil, err
	}

	if err := c.Register(ctx, JwksUrl); err != nil {
		return nil, err
	}

	return c, nil
}

// https://pkg.go.dev/github.com/lestrrat-go/jwx/v3
func UserFromRequest(r *http.Request) (User, error) {
	if JwkCache == nil {
		cache, err := getCachedKeyset(r.Context())
		if err != nil {
			return User{}, fmt.Errorf("jwk cache: %w", err)
		}
		JwkCache = cache
	}

	keyset, err := JwkCache.Lookup(r.Context(), JwksUrl)
	if err != nil {
		return User{}, fmt.Errorf("retrieve jwks: %w", err)
	}

	token, err := jwt.ParseRequest(r, jwt.WithKeySet(keyset))
	if err != nil {
		return User{}, fmt.Errorf("jwt parse request: %w", err)
	}

	userID, exists := token.Subject()
	if !exists {
		return User{}, ErrMissingUserID
	}

	var email, name, role string

	token.Get("email", &email)
	token.Get("name", &name)
	token.Get("role", &role)

	return User{
		ID:    userID,
		Email: email,
		Name:  name,
		Role:  role,
	}, nil
}