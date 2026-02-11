package token

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/artrctx/noliteo-core/internal/database/repository"
	"github.com/artrctx/noliteo-core/internal/jwt"
)

type GenerateTokenResponse struct {
	Ident string `json:"ident"`
	Jwt   string `json:"jwt"`
}

type GenerateTokenRequest struct {
	Token string `json:"token"`
}

// Generate JWT given a token
func (t *TokenService) GenerateTokenHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var reqBody GenerateTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
	}

	token, err := repository.New(t.DB).ValidateToken(req.Context(), reqBody.Token)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed validating token: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	generatedJwt, err := jwt.GenerateToken(jwt.Token{TID: token.ID.String(), Ident: token.Ident.String})

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate jwt token: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(GenerateTokenResponse{token.Ident.String, generatedJwt}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (t *TokenService) ValidateTokenHandler(w http.ResponseWriter, req *http.Request) {
	token, err := jwt.ValidateTokenFromRequest(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed validating jwt: %v", err.Error()), http.StatusUnauthorized)
		return
	}
	if err := json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
