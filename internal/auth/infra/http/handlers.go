// Package http contains HTTP adapters for the auth bounded context: the
// handlers that translate JSON requests into use-case calls, and a Bearer
// middleware that other handlers can use to require authentication.
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/app"
	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

type Handlers struct {
	Register RegisterHandler
	Login    LoginHandler
	Refresh  RefreshHandler
}

type RegisterHandler struct{ UC app.RegisterUser }
type LoginHandler struct{ UC app.LoginUser }
type RefreshHandler struct{ UC app.RefreshSession }

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type tokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
}

func (h RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	var req credentialsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	user, err := h.UC.Execute(r.Context(), app.RegisterInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, userResponse{ID: user.ID, Email: user.Email})
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	var req credentialsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	pair, err := h.UC.Execute(r.Context(), app.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tokenPairToResponse(pair))
}

func (h RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	var req refreshRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	pair, err := h.UC.Execute(r.Context(), req.RefreshToken)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tokenPairToResponse(pair))
}

func tokenPairToResponse(p domain.TokenPair) tokenResponse {
	now := time.Now().UTC()
	return tokenResponse{
		AccessToken:      p.AccessToken,
		RefreshToken:     p.RefreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(p.AccessExpiresAt.Sub(now).Seconds()),
		RefreshExpiresIn: int64(p.RefreshExpiresAt.Sub(now).Seconds()),
	}
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyTaken):
		writeError(w, http.StatusConflict, "email_already_taken")
	case errors.Is(err, domain.ErrWeakPassword):
		writeError(w, http.StatusBadRequest, "weak_password")
	case errors.Is(err, domain.ErrInvalidEmail):
		writeError(w, http.StatusBadRequest, "invalid_email")
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid_credentials")
	case errors.Is(err, domain.ErrInvalidRefresh):
		writeError(w, http.StatusUnauthorized, "invalid_refresh_token")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}
