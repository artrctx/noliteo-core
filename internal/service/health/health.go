package health

import (
	"encoding/json"
	"net/http"

	"github.com/artrctx/noliteo-core/internal/database"
)

func HealthHandlerFunc(s *database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"status":  "healthy",
			"service": "noliteo-core",
			"db":      s.Health(),
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

}
