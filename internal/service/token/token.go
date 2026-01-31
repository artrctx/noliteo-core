package session

import (
	"database/sql"
	"net/http"
)

func VerifyTokenHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

	}
}
