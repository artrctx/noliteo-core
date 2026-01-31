package route

import (
	"encoding/json"
	"net/http"
)

func VerifyAuthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// user, err := auth.UserFromRequest(req)
	// if err != nil {
	// 	slog.Error("failed to get user from request", slog.Any("error", err))
	// 	w.WriteHeader(http.StatusUnauthorized)

	// 	errorRes := response.ErrorResponse{
	// 		Error: "Authentication Failed",
	// 	}

	// 	json.NewEncoder(w).Encode(errorRes)
	// 	return
	// }

	authRes := map[string]string{
		"Status":  "success",
		"Message": "Token is valid",
	}

	json.NewEncoder(w).Encode(authRes)
}
