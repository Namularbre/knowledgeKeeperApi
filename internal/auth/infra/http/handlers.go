// Package http contains HTTP adapters for the auth bounded context: the
// handlers that translate JSON requests into use-case calls, and a Bearer
// middleware that other handlers can use to require authentication.
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	authapp "github.com/Namularbre/knowledgeKeeperApi/internal/auth/app"
	authdomain "github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
	rolesapp "github.com/Namularbre/knowledgeKeeperApi/internal/roles/app"
	rolesdomain "github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type Handlers struct {
	Register   RegisterHandler
	Login      LoginHandler
	Refresh    RefreshHandler
	CreateRole CreateRoleHandler
}

type RegisterHandler struct{ UC authapp.RegisterUser }
type LoginHandler struct{ UC authapp.LoginUser }
type RefreshHandler struct{ UC authapp.RefreshSession }

type CreateRoleHandler struct{ UC rolesapp.CreateRole }

// CredentialsRequest is the body for register/login.
type CredentialsRequest struct {
	Email    string `json:"email" example:"alice@example.com"`
	Password string `json:"password" example:"hunter22"`
}

// RefreshRequest is the body for token refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"vS3...opaque..."`
}

type CreateRoleRequest struct {
	Label string `json:"label" example:"Admin"`
}

// UserResponse is the public view of a user.
type UserResponse struct {
	ID    int64  `json:"id" example:"1"`
	Email string `json:"email" example:"alice@example.com"`
}

// TokenResponse is returned on successful login and refresh.
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type" example:"Bearer"`
	ExpiresIn        int64  `json:"expires_in" example:"900"`
	RefreshExpiresIn int64  `json:"refresh_expires_in" example:"604800"`
}

type CreateRoleResponse struct {
	Id    int64  `json:"id" example:"1"`
	Label string `json:"label" example:"Admin"`
}

// ErrorResponse is the uniform error envelope returned by every handler.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid_credentials"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account. Email must be unique; password must be at least 8 characters.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      CredentialsRequest  true  "Credentials"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request | weak_password | invalid_email"
// @Failure      409   {object}  ErrorResponse  "email_already_taken"
// @Failure      500   {object}  ErrorResponse
// @Router       /auth/register [post]
func (h RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	var req CredentialsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	user, err := h.UC.Execute(r.Context(), authapp.RegisterInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, UserResponse{ID: user.ID, Email: user.Email})
}

// Login godoc
// @Summary      Authenticate a user
// @Description  Exchanges email/password for an access token (short-lived JWT) and an opaque refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      CredentialsRequest  true  "Credentials"
// @Success      200   {object}  TokenResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request"
// @Failure      401   {object}  ErrorResponse  "invalid_credentials"
// @Failure      500   {object}  ErrorResponse
// @Router       /auth/login [post]
func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	var req CredentialsRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	pair, err := h.UC.Execute(r.Context(), authapp.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tokenPairToResponse(pair))
}

// Refresh godoc
// @Summary      Rotate the session tokens
// @Description  Consumes the provided refresh token (revoking it) and returns a fresh access+refresh pair.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RefreshRequest  true  "Refresh token"
// @Success      200   {object}  TokenResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request"
// @Failure      401   {object}  ErrorResponse  "invalid_refresh_token"
// @Failure      500   {object}  ErrorResponse
// @Router       /auth/refresh [post]
func (h RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	var req RefreshRequest
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

func tokenPairToResponse(p authdomain.TokenPair) TokenResponse {
	now := time.Now().UTC()
	return TokenResponse{
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
	writeJSON(w, status, ErrorResponse{Error: code})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authdomain.ErrEmailAlreadyTaken):
		writeError(w, http.StatusConflict, "email_already_taken")
	case errors.Is(err, authdomain.ErrWeakPassword):
		writeError(w, http.StatusBadRequest, "weak_password")
	case errors.Is(err, authdomain.ErrInvalidEmail):
		writeError(w, http.StatusBadRequest, "invalid_email")
	case errors.Is(err, authdomain.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid_credentials")
	case errors.Is(err, authdomain.ErrInvalidRefresh):
		writeError(w, http.StatusUnauthorized, "invalid_refresh_token")
	case errors.Is(err, rolesdomain.ErrRoleAlreadyExists):
		writeError(w, http.StatusConflict, "role_already_exists")
	case errors.Is(err, rolesdomain.ErrInvalidRoleLabel):
		writeError(w, http.StatusBadRequest, "invalid_role_label")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Creates a new role with the provided label. Label is normalized to lowercase and trimmed.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRoleRequest  true  "Role details"
// @Success      200   {object}  CreateRoleResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request | invalid_role_label"
// @Failure      409   {object}  ErrorResponse  "role_already_exists"
// @Failure      500   {object}  ErrorResponse
// @Router       /auth/role/create [post]
func (h CreateRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	var req CreateRoleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	role, err := h.UC.Execute(r.Context(), rolesapp.CreateRoleInput{Label: req.Label})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CreateRoleResponse{Id: role.ID, Label: role.Label})
}
