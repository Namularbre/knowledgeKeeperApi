package domain

import (
	"errors"
	"time"

	"github.com/Namularbre/knowledgeKeeperApi/internal/auth/domain"
)

type DocumentID int64
type DocumentURL string

type document struct {
	ID          DocumentID
	Name        DocumentName
	Description string
	Url         DocumentURL
	Owner       domain.User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	ErrDocumentNotFound      = errors.New("document not found")
	ErrDocumentAlreadyExists = errors.New("document already exists")
	ErrUserDocumentExists    = errors.New("user document already exists")
)
