// Package http contains HTTP adapters for the subjects bounded context.
package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	subjectsapp "github.com/Namularbre/knowledgeKeeperApi/internal/subjects/app"
	subjectsdomain "github.com/Namularbre/knowledgeKeeperApi/internal/subjects/domain"
)

type Handlers struct {
	CreateSubject     CreateSubjectHandler
	FindByID          FindByIDHandler
	FindByUserID      FindByUserIDHandler
	AddUserSubject    AddUserSubjectHandler
	RemoveUserSubject RemoveUserSubjectHandler
	SearchByName      SearchByNameHandler
}

type CreateSubjectHandler struct{ UC subjectsapp.CreateSubject }
type FindByIDHandler struct{ UC subjectsapp.FindByID }
type FindByUserIDHandler struct{ UC subjectsapp.FindByUserID }
type AddUserSubjectHandler struct{ UC subjectsapp.AddUserSubject }
type RemoveUserSubjectHandler struct{ UC subjectsapp.RemoveUserSubject }
type SearchByNameHandler struct{ UC subjectsapp.SearchByName }

type CreateSubjectRequest struct {
	Name string `json:"name" example:"Mathematics"`
}
type SubjectResponse struct {
	Subject subjectsdomain.Subject `json:"subject"`
}
type SubjectsResponse struct {
	Subjects []subjectsdomain.Subject `json:"subjects"`
}
type ErrorResponse struct {
	Error string `json:"error" example:"invalid_request"`
}

// CreateSubject godoc
// @Summary Create a subject
// @Tags subjects
// @Accept json
// @Produce json
// @Param body body CreateSubjectRequest true "Subject details"
// @Success 200 {object} SubjectResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /subjects/create [post]
func (h CreateSubjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var req CreateSubjectRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	subject, err := h.UC.Execute(r.Context(), subjectsapp.CreateSubjectInput{Name: req.Name})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, SubjectResponse{Subject: subject})
}

// FindByID godoc
// @Summary Get a subject by ID
// @Tags subjects
// @Produce json
// @Param id query uint64 true "Subject ID"
// @Success 200 {object} SubjectResponse
// @Failure 404 {object} ErrorResponse
// @Router /subjects/findbyid [get]
func (h FindByIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	id, ok := queryID(w, r, "id")
	if !ok {
		return
	}
	subject, err := h.UC.Execute(r.Context(), subjectsapp.FindByIDInput{ID: id})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, SubjectResponse{Subject: subject})
}

// FindByUserID godoc
// @Summary Get a user's subjects
// @Tags subjects
// @Produce json
// @Param id query uint64 true "User ID"
// @Success 200 {object} SubjectsResponse
// @Router /subjects/findusersubjects [get]
func (h FindByUserIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	id, ok := queryID(w, r, "id")
	if !ok {
		return
	}
	subjects, err := h.UC.Execute(r.Context(), subjectsapp.FindByUserIDInput{ID: id})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, SubjectsResponse{Subjects: subjects})
}

// AddUserSubject godoc
// @Summary Assign a subject to a user
// @Tags subjects
// @Produce json
// @Param user_id query uint64 true "User ID"
// @Param subject_id query uint64 true "Subject ID"
// @Success 200 {object} SubjectsResponse
// @Router /subjects/addusersubject [post]
func (h AddUserSubjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	userID, ok := queryID(w, r, "user_id")
	if !ok {
		return
	}
	subjectID, ok := queryID(w, r, "subject_id")
	if !ok {
		return
	}
	subjects, err := h.UC.Execute(r.Context(), subjectsapp.AddUserSubjectInput{UserID: userID, SubjectID: subjectID})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, SubjectsResponse{Subjects: subjects})
}

// RemoveUserSubject godoc
// @Summary Remove a subject from a user
// @Tags subjects
// @Produce json
// @Param user_id query uint64 true "User ID"
// @Param subject_id query uint64 true "Subject ID"
// @Success 200 {object} SubjectsResponse
// @Router /subjects/removeusersubject [delete]
func (h RemoveUserSubjectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodDelete) {
		return
	}
	userID, ok := queryID(w, r, "user_id")
	if !ok {
		return
	}
	subjectID, ok := queryID(w, r, "subject_id")
	if !ok {
		return
	}
	subjects, err := h.UC.Execute(r.Context(), subjectsapp.RemoveUserSubjectInput{UserID: userID, SubjectID: subjectID})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, SubjectsResponse{Subjects: subjects})
}

// SearchByName godoc
// @Summary Search subjects by name
// @Tags subjects
// @Produce json
// @Param name query string true "Name search term"
// @Success 200 {object} SubjectsResponse
// @Router /subjects/searchbyname [get]
func (h SearchByNameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	subjects, err := h.UC.Execute(r.Context(), subjectsapp.SearchByNameInput{Name: name})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, SubjectsResponse{Subjects: subjects})
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
	case errors.Is(err, subjectsdomain.ErrSubjectNotFound):
		writeError(w, http.StatusNotFound, "subject_not_found")
	case errors.Is(err, subjectsdomain.ErrSubjectAlreadyExists), errors.Is(err, subjectsdomain.ErrUserSubjectExists):
		writeError(w, http.StatusConflict, "subject_already_exists")
	case errors.Is(err, subjectsdomain.ErrInvalidSubjectName):
		writeError(w, http.StatusBadRequest, "invalid_subject_name")
	default:
		log.Printf("subjects: unexpected error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}
