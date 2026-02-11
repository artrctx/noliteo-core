package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/artrctx/noliteo-core/internal/jwt"
)

type TokenKey string

const TokenCtxKey TokenKey = "noliteo-token"

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ValidateTokenFromRequest(r)
		if err != nil {
			slog.Error("failed to get valid token from request", slog.Any("error", err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), TokenCtxKey, token)))
		next.ServeHTTP(w, r)
	})
}
