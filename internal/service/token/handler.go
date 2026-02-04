package token

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/artrctx/quoin-core/internal/database/repository"
	"github.com/artrctx/quoin-core/internal/helper/response"
	"github.com/go-chi/chi/v5"
)

type ValidationResponse struct {
	Ident string `json:"ident"`
	// need to return session token
}

// Validate Token given a token
func (t *TokenService) ValidateTokenHandler(w http.ResponseWriter, req *http.Request) {
	tokenStr := chi.URLParam(req, "token")

	if len(tokenStr) == 0 {
		http.Error(w, "Invalid token length", http.StatusBadRequest)
		return
	}

	token, err := repository.New(t.DB).ValidateToken(req.Context(), tokenStr)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed validating token with error:%v", err.Error()), http.StatusInternalServerError)
		return
	}

	resBody := response.SuccessResponse{
		Message: "Valid token",
		Data:    ValidationResponse{token.Ident.String},
	}

	if err := json.NewEncoder(w).Encode(resBody); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
