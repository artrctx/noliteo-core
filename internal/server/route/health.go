package route

import (
	"encoding/json"
	"net/http"

	"github.com/artrctx/quoin-core/internal/database"
)

func HealthHandlerFunc(s *database.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//TODO: might want to update with
		response := map[string]interface{}{
			"status":  "healthy",
			"service": "quoin-sever",
			"db":      s.Health(),
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

}
