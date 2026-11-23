package domain

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type DocumentName struct {
	value string
}

var ErrInvalidDocumentName = errors.New("invalid document name")

func NewDocumentName(raw string) (DocumentName, error) {
	name := strings.TrimSpace(raw)
	length := utf8.RuneCountInString(name)

	if length == 0 || length > 200 {
		return DocumentName{}, ErrInvalidDocumentName
	}

	return DocumentName{value: name}, nil
}

func (n DocumentName) String() string {
	return n.value
}
