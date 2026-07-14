// Package http contains HTTP adapters for the roles bounded context.
package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	rolesapp "github.com/Namularbre/knowledgeKeeperApi/internal/roles/app"
	rolesdomain "github.com/Namularbre/knowledgeKeeperApi/internal/roles/domain"
)

type Handlers struct {
	CreateRole     CreateRoleHandler
	FindByID       FindByIDHandler
	FindByUserID   FindByUserIDHandler
	AddUserRole    AddUserRoleHandler
	RemoveUserRole RemoveUserRoleHandler
	SearchByLabel  SearchByLabelHandler
}

type CreateRoleHandler struct{ UC rolesapp.CreateRole }

type FindByIDHandler struct {
	UC rolesapp.FindByID
}

type FindByUserIDHandler struct{ UC rolesapp.FindByUserID }

type AddUserRoleHandler struct {
	UC rolesapp.AddUserRole
}

type RemoveUserRoleHandler struct{ UC rolesapp.RemoveUserRole }

type SearchByLabelHandler struct{ UC rolesapp.SearchByLabel }

type CreateRoleRequest struct {
	Label string `json:"label" example:"Admin"`
}

type FindByIDRequest struct {
	ID uint64 `json:"id" example:"1"`
}

type FindByUserIDRequest struct {
	ID uint64 `json:"id" example:"1"`
}

type AddUserRoleRequest struct {
	UserID uint64 `json:"id" example:"1"`
	RoleID uint64 `json:"role_id" example:"1"`
}

type RemoveUserRoleRequest struct {
	RoleID uint64 `json:"role_id" example:"1"`
	UserID uint64 `json:"user_id" example:"1"`
}

type SearchByLabelRequest struct {
	Label string `json:"label" example:"admin"`
}

type CreateRoleResponse struct {
	Id    int64  `json:"id" example:"1"`
	Label string `json:"label" example:"Admin"`
}

type FindByIdResponse struct {
	Role rolesdomain.Role `json:"role"`
}

type FindByUserIDResponse struct {
	Roles []rolesdomain.Role `json:"roles"`
}

type AddUserRoleResponse struct {
	Roles []rolesdomain.Role `json:"roles"`
}

type RemoveUserRoleResponse struct {
	Roles []rolesdomain.Role `json:"roles"`
}

type SearchByLabelResponse struct {
	Roles []rolesdomain.Role `json:"roles"`
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

// FindByID godoc
// @Summary      Get a role by ID
// @Description  Retrieves a single role by its ID.
// @Tags         roles
// @Produce      json
// @Param        id  query      uint64  true  "Role ID"
// @Success      200   {object}  FindByIdResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request"
// @Failure      404   {object}  ErrorResponse  "role_not_found"
// @Failure      500   {object}  ErrorResponse
// @Router       /roles/findbyid [get]
func (h FindByIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	id := uint64(0)
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	role, err := h.UC.Execute(r.Context(), rolesapp.FindByIDInput{ID: id})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, FindByIdResponse{Role: role})
}

// FindByUserID godoc
// @Summary      Get the roles of a user
// @Description  Retrieves all roles assigned to a given user.
// @Tags         roles
// @Produce      json
// @Param        id  query      uint64  true  "User ID"
// @Success      200   {object}  FindByUserIDResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request"
// @Failure      500   {object}  ErrorResponse
// @Router       /roles/finduserroles [get]
func (h FindByUserIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	id := uint64(0)
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	roles, err := h.UC.Execute(r.Context(), rolesapp.FindByUserIDInput{ID: id})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, FindByUserIDResponse{Roles: roles})
}

// AddUserRole godoc
// @Summary      Assign a role to a user
// @Description  Adds a role to a user and returns the user's updated list of roles.
// @Tags         roles
// @Produce      json
// @Param        user_id  query      uint64  true  "User ID"
// @Param        role_id  query      uint64  true  "Role ID"
// @Success      200   {object}  AddUserRoleResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request"
// @Failure      500   {object}  ErrorResponse
// @Router       /roles/adduserrole [post]
func (h AddUserRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	userIdStr := r.URL.Query().Get("user_id")
	roleIdStr := r.URL.Query().Get("role_id")
	if userIdStr == "" || roleIdStr == "" {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	userId, roleId := uint64(0), uint64(0)
	if _, err := fmt.Sscanf(userIdStr, "%d", &userId); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	if _, err := fmt.Sscanf(roleIdStr, "%d", &roleId); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	roles, err := h.UC.Execute(r.Context(), rolesapp.AddUserRoleInput{UserID: userId, RoleID: roleId})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, AddUserRoleResponse{Roles: roles})
}

// RemoveUserRole godoc
// @Summary      Remove a role from a user
// @Description  Removes a role from a user and returns the user's updated list of roles.
// @Tags         roles
// @Produce      json
// @Param        user_id  query      uint64  true  "User ID"
// @Param        role_id  query      uint64  true  "Role ID"
// @Success      200   {object}  RemoveUserRoleResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request"
// @Failure      500   {object}  ErrorResponse
// @Router       /roles/removeuserrole [delete]
func (h RemoveUserRoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	userIdStr := r.URL.Query().Get("user_id")
	roleIdStr := r.URL.Query().Get("role_id")
	if userIdStr == "" || roleIdStr == "" {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	userId, roleId := uint64(0), uint64(0)
	if _, err := fmt.Sscanf(userIdStr, "%d", &userId); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	if _, err := fmt.Sscanf(roleIdStr, "%d", &roleId); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	roles, err := h.UC.Execute(r.Context(), rolesapp.RemoveUserRoleInput{UserID: userId, RoleID: roleId})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, RemoveUserRoleResponse{Roles: roles})
}

// SearchByLabel godoc
// @Summary      Search roles by label
// @Description  Retrieves roles whose label contains the given substring.
// @Tags         roles
// @Produce      json
// @Param        label  query      string  true  "Label search term"
// @Success      200   {object}  SearchByLabelResponse
// @Failure      400   {object}  ErrorResponse  "invalid_request"
// @Failure      500   {object}  ErrorResponse
// @Router       /roles/searchbylabel [get]
func (h SearchByLabelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}
	label := r.URL.Query().Get("label")
	if label == "" {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	roles, err := h.UC.Execute(r.Context(), rolesapp.SearchByLabelInput{Label: label})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, SearchByLabelResponse{Roles: roles})
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
		log.Printf("roles: unexpected error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}
