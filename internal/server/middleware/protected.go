package middleware

import (
	"net/http"
)

type UserKey string

const UserCtxKey UserKey = "user"

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// user, err := auth.UserFromRequest(r)

		// if err != nil {
		// 	slog.Error("failed to get user from request", slog.Any("error", err))
		// 	w.WriteHeader(http.StatusUnauthorized)

		// 	errorRes := response.ErrorResponse{
		// 		Error: "Authentication Failed",
		// 	}

		// 	json.NewEncoder(w).Encode(errorRes)
		// 	return
		// }

		// next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserCtxKey, user)))
		next.ServeHTTP(w, r)
	})
}
