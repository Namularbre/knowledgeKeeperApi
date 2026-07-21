// Package http contains HTTP adapters for the cohorts bounded context.
package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	cohortapp "github.com/Namularbre/knowledgeKeeperApi/internal/cohort/app"
	cohortdomain "github.com/Namularbre/knowledgeKeeperApi/internal/cohort/domain"
)

type Handlers struct {
	CreateCohort     CreateCohortHandler
	FindByID         FindByIDHandler
	FindByUserID     FindByUserIDHandler
	AddUserCohort    AddUserCohortHandler
	RemoveUserCohort RemoveUserCohortHandler
	SearchByName     SearchByNameHandler
}

type CreateCohortHandler struct{ UC cohortapp.CreateCohort }
type FindByIDHandler struct{ UC cohortapp.FindByID }
type FindByUserIDHandler struct{ UC cohortapp.FindByUserID }
type AddUserCohortHandler struct{ UC cohortapp.AddUserCohort }
type RemoveUserCohortHandler struct{ UC cohortapp.RemoveUserCohort }
type SearchByNameHandler struct{ UC cohortapp.SearchByName }

type CreateCohortRequest struct {
	Name string `json:"name" example:"First year"`
}
type CohortResponse struct {
	Cohort cohortdomain.Cohort `json:"cohort"`
}
type CohortsResponse struct {
	Cohorts []cohortdomain.Cohort `json:"cohorts"`
}
type ErrorResponse struct {
	Error string `json:"error" example:"invalid_request"`
}

// CreateCohort godoc
// @Summary Create a cohort
// @Tags cohorts
// @Accept json
// @Produce json
// @Param body body CreateCohortRequest true "Cohort details"
// @Success 200 {object} CohortResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /cohorts/create [post]
func (h CreateCohortHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var req CreateCohortRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	cohort, err := h.UC.Execute(r.Context(), cohortapp.CreateCohortInput{Name: req.Name})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CohortResponse{Cohort: cohort})
}

// FindByID godoc
// @Summary Get a cohort by ID
// @Tags cohorts
// @Produce json
// @Param id query uint64 true "Cohort ID"
// @Success 200 {object} CohortResponse
// @Failure 404 {object} ErrorResponse
// @Router /cohorts/findbyid [get]
func (h FindByIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	id, ok := queryID(w, r, "id")
	if !ok {
		return
	}
	cohort, err := h.UC.Execute(r.Context(), cohortapp.FindByIDInput{ID: id})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CohortResponse{Cohort: cohort})
}

// FindByUserID godoc
// @Summary Get a user's cohorts
// @Tags cohorts
// @Produce json
// @Param id query uint64 true "User ID"
// @Success 200 {object} CohortsResponse
// @Router /cohorts/findusercohorts [get]
func (h FindByUserIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	id, ok := queryID(w, r, "id")
	if !ok {
		return
	}
	cohorts, err := h.UC.Execute(r.Context(), cohortapp.FindByUserIDInput{ID: id})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CohortsResponse{Cohorts: cohorts})
}

// AddUserCohort godoc
// @Summary Assign a cohort to a user
// @Tags cohorts
// @Produce json
// @Param user_id query uint64 true "User ID"
// @Param cohort_id query uint64 true "Cohort ID"
// @Success 200 {object} CohortsResponse
// @Router /cohorts/addusercohort [post]
func (h AddUserCohortHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	userID, ok := queryID(w, r, "user_id")
	if !ok {
		return
	}
	cohortID, ok := queryID(w, r, "cohort_id")
	if !ok {
		return
	}
	cohorts, err := h.UC.Execute(r.Context(), cohortapp.AddUserCohortInput{UserID: userID, CohortID: cohortID})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CohortsResponse{Cohorts: cohorts})
}

// RemoveUserCohort godoc
// @Summary Remove a cohort from a user
// @Tags cohorts
// @Produce json
// @Param user_id query uint64 true "User ID"
// @Param cohort_id query uint64 true "Cohort ID"
// @Success 200 {object} CohortsResponse
// @Router /cohorts/removeusercohort [delete]
func (h RemoveUserCohortHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodDelete) {
		return
	}
	userID, ok := queryID(w, r, "user_id")
	if !ok {
		return
	}
	cohortID, ok := queryID(w, r, "cohort_id")
	if !ok {
		return
	}
	cohorts, err := h.UC.Execute(r.Context(), cohortapp.RemoveUserCohortInput{UserID: userID, CohortID: cohortID})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CohortsResponse{Cohorts: cohorts})
}

// SearchByName godoc
// @Summary Search cohorts by name
// @Tags cohorts
// @Produce json
// @Param name query string true "Name search term"
// @Success 200 {object} CohortsResponse
// @Router /cohorts/searchbyname [get]
func (h SearchByNameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	cohorts, err := h.UC.Execute(r.Context(), cohortapp.SearchByNameInput{Name: name})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CohortsResponse{Cohorts: cohorts})
}

func queryID(w http.ResponseWriter, r *http.Request, key string) (uint64, bool) {
	id, err := strconv.ParseUint(r.URL.Query().Get(key), 10, 64)
	if err != nil || id == 0 {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return 0, false
	}
	return id, true
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return false
	}
	return true
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
	case errors.Is(err, cohortdomain.ErrCohortNotFound):
		writeError(w, http.StatusNotFound, "cohort_not_found")
	case errors.Is(err, cohortdomain.ErrCohortAlreadyExists), errors.Is(err, cohortdomain.ErrUserCohortExists):
		writeError(w, http.StatusConflict, "cohort_already_exists")
	case errors.Is(err, cohortdomain.ErrInvalidCohortName):
		writeError(w, http.StatusBadRequest, "invalid_cohort_name")
	default:
		log.Printf("cohorts: unexpected error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}
