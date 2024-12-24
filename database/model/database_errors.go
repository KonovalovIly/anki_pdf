package database_models

import (
	"database/sql"
	"strings"
)

type DatabaseError struct {
	Typ        string
	Error      string
	Invocation string
}

func ProcessErrorFromDatabase(err error, invocation string) *DatabaseError {
	switch {
	case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
		return &DatabaseError{Typ: "duplicate_key", Error: err.Error(), Invocation:  invocation}
	case err == sql.ErrNoRows:
		return &DatabaseError{Typ: "no_row", Error: err.Error(), Invocation:  invocation}
	default:
		return &DatabaseError{Typ: "unknown", Error: err.Error(), Invocation:  invocation}
	}
}
