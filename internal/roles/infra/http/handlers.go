// Package http contains HTTP adapters for the roles bounded context.
package http

import (
	"encoding/json"
	"errors"
	"net/http"

	rolesapp "github.com/Namularbre/knowledgeKeeperApi/internal/roles/app"
	rolesdomain "github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type Handlers struct {
	CreateRole CreateRoleHandler
}

type CreateRoleHandler struct{ UC rolesapp.CreateRole }

type CreateRoleRequest struct {
	Label string `json:"label" example:"Admin"`
}

type CreateRoleResponse struct {
	Id    int64  `json:"id" example:"1"`
	Label string `json:"label" example:"Admin"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid_credentials"`
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Creates a new role with the provided label. Label is normalized to lowercase and trimmed.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRoleRequest  true  "Role details"
// @Success      200   {object}  CreateRoleResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request | invalid_role_label"
// @Failure      409   {object}  ErrorResponse  "role_already_exists"
// @Failure      500   {object}  ErrorResponse
// @Router       /roles/create [post]
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
	case errors.Is(err, rolesdomain.ErrRoleAlreadyExists):
		writeError(w, http.StatusConflict, "role_already_exists")
	case errors.Is(err, rolesdomain.ErrInvalidRoleLabel):
		writeError(w, http.StatusBadRequest, "invalid_role_label")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}
