package model

import (
	"database/sql"
	"strings"
)

type DatabaseError struct {
	Typ   string
	Error string
}

func ProcessErrorFromDatabase(err error) *DatabaseError {
	switch {
	case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
		return &DatabaseError{Typ: "duplicate_key", Error: err.Error()}
	case err == sql.ErrNoRows:
		return &DatabaseError{Typ: "no_row", Error: err.Error()}
	default:
		return &DatabaseError{Typ: "unknown", Error: err.Error()}
	}
}
